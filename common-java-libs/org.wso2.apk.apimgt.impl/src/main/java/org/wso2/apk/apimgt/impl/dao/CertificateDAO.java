/*
 * Copyright (c) 2022, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.apimgt.impl.dao;

import org.wso2.apk.apimgt.api.dto.CertificateMetadataDTO;
import org.wso2.apk.apimgt.api.dto.ClientCertificateDTO;
import org.wso2.apk.apimgt.api.model.Identifier;
import org.wso2.apk.apimgt.impl.certificatemgt.exceptions.CertificateAliasExistsException;
import org.wso2.apk.apimgt.impl.certificatemgt.exceptions.CertificateManagementException;

import java.util.List;

public interface CertificateDAO {

    /**
     * Method to add a new client certificate to the database.
     *
     * @param certificate   : Client certificate that need to be added.
     * @param apiIdentifier : API which the client certificate is uploaded against.
     * @param alias         : Alias for the new certificate.
     * @param tenantId      : The Id of the tenant who uploaded the certificate.
     * @param organization  : Organization
     * @return : True if the information is added successfully, false otherwise.
     * @throws CertificateManagementException if existing entry is found for the given endpoint or alias.
     */
    boolean addClientCertificate(String certificate, Identifier apiIdentifier, String alias, String tierName,
                                 int tenantId, String organization) throws CertificateManagementException;

    /**
     * Method to add a new certificate to the database.
     *
     * @param alias    : Alias for the new certificate.
     * @param endpoint : The endpoint/ server url which the certificate will be mapped to.
     * @param tenantId : The Id of the tenant who uploaded the certificate.
     * @return : True if the information is added successfully, false otherwise.
     * @throws CertificateManagementException if existing entry is found for the given endpoint or alias.
     */
    boolean addCertificate(String certificate, String alias, String endpoint, int tenantId)
            throws CertificateManagementException,
            CertificateAliasExistsException;

    /**
     * Method to retrieve certificate metadata from db for specific tenant which matches alias or api identifier.
     * Both alias and api identifier are optional
     *
     * @param tenantId      : The id of the tenant which the certificate belongs to.
     * @param alias         : Alias for the certificate. (Optional)
     * @param apiIdentifier : The API which the certificate is mapped to. (Optional)
     * @param organization  : Organization
     * @return : A CertificateMetadataDTO object if the certificate is retrieved successfully, null otherwise.
     */
    List<ClientCertificateDTO> getClientCertificates(int tenantId, String alias, Identifier apiIdentifier,
                                                     String organization) throws CertificateManagementException;

    /**
     * Method to retrieve certificate metadata from db for specific tenant which matches alias or endpoint.
     * From alias and endpoint, only one parameter is required.
     *
     * @param tenantId : The id of the tenant which the certificate belongs to.
     * @param alias    : Alias for the certificate. (Optional)
     * @param endpoint : The endpoint/ server url which the certificate is mapped to. (Optional)
     * @return : A CertificateMetadataDTO object if the certificate is retrieved successfully, null otherwise.
     */
    List<CertificateMetadataDTO> getCertificates(String alias, String endpoint, int tenantId)
            throws CertificateManagementException;

    /**
     * To update an already existing client certificate.
     *
     * @param certificate : Specific certificate.
     * @param alias       : Alias of the certificate.
     * @param tier        : Name of tier related with the certificate.
     * @param tenantId    : ID of the tenant.
     * @param organization : Organization
     * @return true if the update succeeds, unless false.
     * @throws CertificateManagementException Certificate Management Exception.
     */
    boolean updateClientCertificate(String certificate, String alias, String tier, int tenantId,
                                    String organization) throws CertificateManagementException;
}
