package org.wso2.apk.integration.utils;

import java.io.File;

public class MultipartFilePart {

    private String name;
    private File file;

    public MultipartFilePart(String name, File file) {

        this.name = name;
        this.file = file;
    }

    public String getName() {

        return name;
    }

    public File getFile() {

        return file;
    }
}
