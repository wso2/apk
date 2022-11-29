package org.wso2.apk.apimgt.impl.dao;

import org.wso2.apk.apimgt.api.APIManagementException;
import org.wso2.apk.apimgt.api.model.Identifier;
import org.wso2.apk.apimgt.api.model.Workflow;
import org.wso2.apk.apimgt.impl.dto.WorkflowDTO;

public interface WorkflowDAO {

    /**
     * Returns a workflow object for a given internal workflow reference and the workflow type.
     *
     * @param workflowReference
     * @param workflowType
     * @return
     * @throws APIManagementException
     */
    WorkflowDTO retrieveWorkflowFromInternalReference(String workflowReference, String workflowType)
            throws APIManagementException;

    /**
     * Retries the WorkflowExternalReference for a subscription.
     *
     * @param subscriptionId ID of the subscription
     * @return External workflow reference for the subscription <code>subscriptionId</code>
     * @throws APIManagementException
     */
    String getExternalWorkflowReferenceForSubscription(int subscriptionId) throws APIManagementException;

    /**
     * Get the Pending workflow Requests using WorkflowType for a particular tenant
     *
     * @param workflowType Type of the workflow pending request
     * @param status       workflow status of workflow pending request
     * @param tenantDomain tenantDomain of the user
     * @return List of workflow pending request
     * @throws APIManagementException
     */
    Workflow[] getWorkflows(String workflowType, String status, String tenantDomain) throws APIManagementException;

    /**
     * Get the Pending workflow Request using ExternalWorkflowReference for a particular tenant
     *
     * @param externalWorkflowRef of pending workflow request
     * @param status              workflow status of workflow pending process
     * @param tenantDomain        tenant domain of user
     * @return workflow pending request
     */
    Workflow getWorkflowReferenceByExternalWorkflowReferenceID(String externalWorkflowRef, String status,
                                                               String tenantDomain) throws APIManagementException;

    /**
     * Retries the WorkflowExternalReference for a subscription.
     *
     * @param identifier Identifier to find the subscribed api
     * @param appID      ID of the application which has the subscription
     * @param organization organization
     * @return External workflow reference for the subscription identified
     * @throws APIManagementException
     */
    String getExternalWorkflowReferenceForSubscription(Identifier identifier, int appID, String organization)
            throws APIManagementException;

    /**
     * Returns a workflow object for a given external workflow reference.
     *
     * @param workflowReference
     * @return
     * @throws APIManagementException
     */
    WorkflowDTO retrieveWorkflow(String workflowReference) throws APIManagementException;

    /**
     * Get external workflow reference by internal workflow reference and workflow type
     * @param internalRef Internal reference of the workflow
     * @param workflowType Workflow type of the workflow
     * @return External workflow reference for the given internal reference and workflow type if present. Null otherwise
     * @throws APIManagementException If an SQL exception occurs in database interactions
     */
    String getExternalWorkflowRefByInternalRefWorkflowType(int internalRef, String workflowType) throws APIManagementException;

    /**
     * Get external workflow reference by Subscription and workflow type
     * @param subscriptionId Subscription ID
     * @param wfType Workflow type of the workflow
     * @return External workflow reference for the given Subscription and workflow type if present. Null otherwise
     * @throws APIManagementException If an SQL exception occurs in database interactions
     */
    String getExternalWorkflowReferenceForSubscriptionAndWFType(int subscriptionId, String wfType) throws APIManagementException;

    /**
     * Retrieves registration workflow reference for applicationId and key type
     *
     * @param applicationId  id of the application with registration
     * @param keyType        key type of the registration
     * @param keyManagerName
     * @return workflow reference of the registration
     * @throws APIManagementException
     */
    String getRegistrationWFReference(int applicationId, String keyType, String keyManagerName)
            throws APIManagementException ;


}
