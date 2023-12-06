/*
 * Copyright (c) 2021, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.enforcer.security.jwt.validator;

import com.nimbusds.jose.JOSEException;
import com.nimbusds.jose.jwk.ECKey;
import com.nimbusds.jose.jwk.JWK;
import com.nimbusds.jose.jwk.JWKSet;
import com.nimbusds.jose.jwk.RSAKey;
import com.nimbusds.jwt.JWTClaimsSet;
import com.nimbusds.jwt.SignedJWT;
import com.nimbusds.jwt.util.DateUtils;
import org.apache.commons.lang3.StringUtils;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.wso2.apk.enforcer.commons.dto.JWTValidationInfo;
import org.wso2.apk.enforcer.commons.exception.JWTGeneratorException;
import org.wso2.apk.enforcer.commons.jwttransformer.JWTTransformer;
import org.wso2.apk.enforcer.commons.exception.EnforcerException;
import org.wso2.apk.enforcer.config.ConfigHolder;
import org.wso2.apk.enforcer.config.dto.ExtendedTokenIssuerDto;
import org.wso2.apk.enforcer.constants.APIConstants;
import org.wso2.apk.enforcer.constants.APISecurityConstants;
import org.wso2.apk.enforcer.security.jwt.SignedJWTInfo;
import org.wso2.apk.enforcer.util.JWKSClient;
import org.wso2.apk.enforcer.util.JWTUtils;

import java.security.PublicKey;
import java.security.cert.Certificate;
import java.security.interfaces.ECPublicKey;
import java.security.interfaces.RSAPublicKey;
import java.text.ParseException;
import java.util.Arrays;
import java.util.Collections;
import java.util.Date;
import java.util.List;

/**
 * Class responsible to validate jwt. This should validate the JWT signature, expiry time.
 * validating the sub, aud, iss can be made optional.
 */
public class JWTValidator {
    private static final Logger logger = LogManager.getLogger(JWTValidator.class);
    private JWKSet jwkSet;
    JWTTransformer jwtTransformer;
    ExtendedTokenIssuerDto tokenIssuer;
    JWKSClient jwksClient;

    public JWTValidator(ExtendedTokenIssuerDto tokenIssuer) throws EnforcerException {
        jwtTransformer = ConfigHolder.getInstance().getConfig().getJwtTransformer(tokenIssuer.getIssuer());
        jwtTransformer.loadConfiguration(tokenIssuer);
        this.tokenIssuer = tokenIssuer;
        if (tokenIssuer.getJwksConfigurationDTO() != null && tokenIssuer.getJwksConfigurationDTO().isEnabled() && StringUtils.isNotEmpty(tokenIssuer.getJwksConfigurationDTO().getUrl())) {
            Certificate certificate = tokenIssuer.getJwksConfigurationDTO().getCertificate();
            if (certificate != null) {
                jwksClient = new JWKSClient(tokenIssuer.getJwksConfigurationDTO().getUrl(), List.of(certificate));
            } else {
                jwksClient = new JWKSClient(tokenIssuer.getJwksConfigurationDTO().getUrl(), Collections.emptyList());
            }
        }
    }

    public JWTValidationInfo validateToken(String token, SignedJWTInfo signedJWTInfo) throws EnforcerException {
        JWTValidationInfo jwtValidationInfo = new JWTValidationInfo();
        boolean state;
        try {
            state = validateSignature(signedJWTInfo.getSignedJWT());
            if (state) {
                JWTClaimsSet jwtClaimsSet = signedJWTInfo.getJwtClaimsSet();
                state = validateTokenExpiry(jwtClaimsSet);
                if (state) {
                    jwtValidationInfo.setConsumerKey(jwtTransformer.getTransformedConsumerKey(jwtClaimsSet));
                    jwtValidationInfo.setScopes(jwtTransformer.getTransformedScopes(jwtClaimsSet));
                    JWTClaimsSet transformedJWTClaimSet = jwtTransformer.transform(jwtClaimsSet);
                    createJWTValidationInfoFromJWT(jwtValidationInfo, transformedJWTClaimSet);
                    jwtValidationInfo.setKeyManager(tokenIssuer.getName());
                    jwtValidationInfo.setIdentifier(JWTUtils.getJWTTokenIdentifier(signedJWTInfo));
                    jwtValidationInfo.setJwtClaimsSet(signedJWTInfo.getJwtClaimsSet());
                    jwtValidationInfo.setToken(token);
                    return jwtValidationInfo;
                }
                jwtValidationInfo.setValidationCode(APISecurityConstants.API_AUTH_ACCESS_TOKEN_EXPIRED);
                logger.debug("Token is expired.");
            } else {
                jwtValidationInfo.setValidationCode(APIConstants.KeyValidationStatus.API_AUTH_INVALID_CREDENTIALS);
                logger.debug("Token signature is invalid.");
            }
        } catch (ParseException | JWTGeneratorException e) {
            throw new EnforcerException("Error while parsing JWT", e);
        }
        jwtValidationInfo.setValid(false);
        return jwtValidationInfo;
    }

    protected boolean validateSignature(SignedJWT signedJWT) throws EnforcerException {
        try {
            String keyID = signedJWT.getHeader().getKeyID();
            if (jwksClient != null) {
                // Check JWKSet Available in Cache
                if (jwkSet == null) {
                    jwkSet = jwksClient.getJWKSet();
                }
                JWK jwkSetKeyByKeyId = jwkSet.getKeyByKeyId(keyID);
                if (jwkSetKeyByKeyId == null) {
                    jwkSet = jwksClient.getJWKSet();
                }
                jwkSetKeyByKeyId = jwkSet.getKeyByKeyId(keyID);
                if (jwkSetKeyByKeyId instanceof RSAKey) {
                    RSAKey keyByKeyId = (RSAKey) jwkSetKeyByKeyId;
                    RSAPublicKey rsaPublicKey = keyByKeyId.toRSAPublicKey();
                    if (rsaPublicKey != null) {
                        return JWTUtils.verifyTokenSignature(signedJWT, rsaPublicKey);
                    }
                } else if (jwkSetKeyByKeyId instanceof ECKey) {
                    ECKey keyByKeyId = (ECKey) jwkSetKeyByKeyId;
                    ECPublicKey ecPublicKey = keyByKeyId.toECPublicKey();
                    if (ecPublicKey != null) {
                        return JWTUtils.verifyTokenSignature(signedJWT, ecPublicKey);
                    }
                } else {
                    throw new EnforcerException("Key Algorithm not supported");
                }
            }
            if (tokenIssuer.getCertificate() != null) {
                logger.debug("Retrieve certificate from Token issuer and validating");
                PublicKey publicKey = tokenIssuer.getCertificate().getPublicKey();
                return JWTUtils.verifyTokenSignature(signedJWT, publicKey);
            } else {
                throw new EnforcerException("Certificate not found for validation");
            }
        } catch (JOSEException e) {
            throw new EnforcerException("JWT Signature verification failed", e);
        }
    }

    protected boolean validateTokenExpiry(JWTClaimsSet jwtClaimsSet) {

        long timestampSkew = 5; //TODO : Read from config.
        Date now = new Date();
        Date exp = jwtClaimsSet.getExpirationTime();
        return exp == null || DateUtils.isAfter(exp, now, timestampSkew);
    }

    private void createJWTValidationInfoFromJWT(JWTValidationInfo jwtValidationInfo, JWTClaimsSet jwtClaimsSet)
            throws ParseException {
        jwtValidationInfo.setValid(true);
        jwtValidationInfo.setClaims(jwtClaimsSet.getClaims());
        jwtValidationInfo.setExpiryTime(jwtClaimsSet.getExpirationTime().getTime());
        jwtValidationInfo.setUser(jwtClaimsSet.getSubject());
        if (jwtClaimsSet.getClaim(APIConstants.JwtTokenConstants.SCOPE) != null) {
            if (jwtClaimsSet.getClaim(APIConstants.JwtTokenConstants.SCOPE) instanceof List) {
                jwtValidationInfo.setScopes(jwtClaimsSet.getStringListClaim(APIConstants.JwtTokenConstants.SCOPE));
            } else {
                jwtValidationInfo.setScopes(Arrays.asList(
                        jwtClaimsSet.getStringClaim(APIConstants.JwtTokenConstants.SCOPE)
                                .split(APIConstants.JwtTokenConstants.SCOPE_DELIMITER)));
            }
        }
    }


}
