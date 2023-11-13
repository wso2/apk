package org.wso2.apk.enforcer.subscription;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;

public class ApplicationListDto implements Serializable {
private List<ApplicationDto> list = new ArrayList<>();

    public List<ApplicationDto> getList() {

        return list;
    }

    public void setList(List<ApplicationDto> list) {

        this.list = list;
    }
}
