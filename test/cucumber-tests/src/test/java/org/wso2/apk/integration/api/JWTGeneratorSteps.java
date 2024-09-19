package org.wso2.apk.integration.api;

import com.google.common.io.Resources;
import com.nimbusds.jose.JOSEException;
import com.nimbusds.jose.JWSAlgorithm;
import com.nimbusds.jose.JWSHeader;
import com.nimbusds.jose.JWSSigner;
import com.nimbusds.jose.crypto.RSASSASigner;
import com.nimbusds.jose.jwk.RSAKey;
import com.nimbusds.jwt.JWTClaimsSet;
import com.nimbusds.jwt.SignedJWT;
import io.cucumber.java.en.And;
import io.cucumber.java.en.Then;
import io.cucumber.java.en.When;

import java.io.File;
import java.io.IOException;
import java.net.URL;
import java.security.KeyStore;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;
import java.security.cert.CertificateException;
import java.util.Date;
import java.util.UUID;

import org.wso2.apk.integration.utils.Constants;

public class JWTGeneratorSteps {

    private final SharedContext sharedContext;

    public JWTGeneratorSteps(SharedContext sharedContext) {

        this.sharedContext = sharedContext;
    }

    @Then("I generate JWT token from idp1 with kid {string}")
    public void generateTokenFromIdp1(String kid) throws IOException, CertificateException, KeyStoreException,
            NoSuchAlgorithmException, JOSEException {

        URL url = Resources.getResource("artifacts/jwtcert/idp1.jks");
        File keyStoreFile = new File(url.getPath());
        KeyStore keyStore = KeyStore.getInstance(keyStoreFile, "wso2carbon".toCharArray());
        RSAKey rsaKey = RSAKey.load(keyStore, "idp1Key", "wso2carbon".toCharArray());
        JWSSigner signer = new RSASSASigner(rsaKey);
        JWTClaimsSet claimsSet = new JWTClaimsSet.Builder()
                .subject("alice")
                .issuer("https://idp1.com")
                .expirationTime(new Date(new Date().getTime() + 60 * 1000))
                .jwtID(UUID.randomUUID().toString())
                .claim("azp", UUID.randomUUID().toString())
                .claim("scope", Constants.API_CREATE_SCOPE)
                .build();
        SignedJWT signedJWT = new SignedJWT(
                new JWSHeader.Builder(JWSAlgorithm.RS256).keyID(kid).build(),
                claimsSet);
        signedJWT.sign(signer);
        String jwtToken = signedJWT.serialize();
        sharedContext.addStoreValue("idp-1-token", jwtToken);
    }
    @Then("I generate JWT token from idp1 with kid {string} and consumer_key {string}")
    public void generateTokenFromIdp1WithConsumerKey(String kid,String consumerKey) throws IOException, CertificateException, KeyStoreException,
            NoSuchAlgorithmException, JOSEException {

        URL url = Resources.getResource("artifacts/jwtcert/idp1.jks");
        File keyStoreFile = new File(url.getPath());
        KeyStore keyStore = KeyStore.getInstance(keyStoreFile, "wso2carbon".toCharArray());
        RSAKey rsaKey = RSAKey.load(keyStore, "idp1Key", "wso2carbon".toCharArray());
        JWSSigner signer = new RSASSASigner(rsaKey);
        JWTClaimsSet claimsSet = new JWTClaimsSet.Builder()
                .subject("alice")
                .issuer("https://idp1.com")
                .expirationTime(new Date(new Date().getTime() + 60 * 60 * 24  * 1000))
                .jwtID(UUID.randomUUID().toString())
                .claim("azp", consumerKey)
                .claim("scope", Constants.API_CREATE_SCOPE)
                .build();
        SignedJWT signedJWT = new SignedJWT(
                new JWSHeader.Builder(JWSAlgorithm.RS256).keyID(kid).build(),
                claimsSet);
        signedJWT.sign(signer);
        String jwtToken = signedJWT.serialize();
        sharedContext.addStoreValue("idp-1-"+consumerKey+"-token", jwtToken);
    }

    public static void main(String[] args) throws CertificateException, IOException, KeyStoreException, NoSuchAlgorithmException, JOSEException {
        SharedContext sharedContext1 = new SharedContext();
        String consumerKey = "45f1c5c8-a92e-11ed-afa1-0242ac120005";
        new JWTGeneratorSteps(sharedContext1).generateTokenFromIdp1WithConsumerKey("123-456", consumerKey);
        System.out.println(sharedContext1.getStoreValue("idp-1-"+consumerKey+"-token"));
    }


    @And("I have a valid token for organization {string}")
    public void generateTokenFromIdp1WithOrganization(String organization) throws IOException, CertificateException,
            KeyStoreException,
            NoSuchAlgorithmException, JOSEException {

        URL url = Resources.getResource("artifacts/jwtcert/idp1.jks");
        File keyStoreFile = new File(url.getPath());
        KeyStore keyStore = KeyStore.getInstance(keyStoreFile, "wso2carbon".toCharArray());
        RSAKey rsaKey = RSAKey.load(keyStore, "idp1Key", "wso2carbon".toCharArray());
        JWSSigner signer = new RSASSASigner(rsaKey);
        JWTClaimsSet claimsSet = new JWTClaimsSet.Builder()
                .subject("alice")
                .issuer("https://idp1.com")
                .expirationTime(new Date(new Date().getTime() + 60 * 1000))
                .jwtID(UUID.randomUUID().toString())
                .claim("azp", UUID.randomUUID().toString())
                .claim("scope", Constants.API_CREATE_SCOPE)
                .claim("organization", organization)
                .build();
        SignedJWT signedJWT = new SignedJWT(
                new JWSHeader.Builder(JWSAlgorithm.RS256).keyID("123-456").build(),
                claimsSet);
        signedJWT.sign(signer);
        String jwtToken = signedJWT.serialize();
        sharedContext.addStoreValue(organization, jwtToken);
    }
}
