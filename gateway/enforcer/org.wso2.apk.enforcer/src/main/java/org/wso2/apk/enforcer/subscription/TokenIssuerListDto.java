package org.wso2.apk.enforcer.subscription;

import org.wso2.apk.enforcer.commons.dto.TokenIssuerDto;
import org.wso2.apk.enforcer.discovery.subscription.TokenIssuer;

import java.util.List;

public class TokenIssuerListDto {
private List<TokenIssuerRestDto> list;

    public List<TokenIssuerRestDto> getList() {

        return list;
    }
}
