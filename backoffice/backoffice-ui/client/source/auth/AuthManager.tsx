
import LoggedInUser from 'types/LoggedInUser';
export const getUser = (): LoggedInUser => {
    return {
        name: 'admin@carbon.super',
        _scopes: ['apim:admin', 'apim:admin_alert_manage', 'apim:admin_application_view',
            'apim:admin_operations', 'apim:admin_settings', 'apim:api_import_export',
            'apim:api_product_import_export', 'apim:api_workflow_approve',
            'apim:api_workflow_view', 'apim:app_import_export', 'apim:app_owner_change',
            'apim:bl_manage', 'apim:bl_view', 'apim:bot_data',
            'apim:environment_manage', 'apim:environment_read', 'apim:mediation_policy_create',
            'apim:mediation_policy_view',
            'apim:monetization_usage_publish', 'apim:policies_import_export', 'apim:role_manage',
            'apim:scope_manage', 'apim:tenantInfo',
            'apim:tenant_theme_manage', 'apim:tier_manage', 'apim:tier_view', 'openid'],
        _remember: false,
        _environmentName: 'Default',
        rememberMe: false,
    };
}

