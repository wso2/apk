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

package org.wso2.apk.apimgt.api.doc.model;

import java.util.List;
import java.util.Objects;

public class APIResource {
	
	private String path;
	
	private String description;
	
	private List<Operation> operations;

	private String verb;
	
	public APIResource(String path, String description, List<Operation> ops) {
		this.path = path;
		this.description = description;
		this.operations = ops;
	}

	public APIResource(String verb, String path) {
		this.verb = verb;
		this.path = path;
	}

	@Override
	public String toString() {
		return "{" +
				"verb='" + verb + '\'' +
				", path='" + path + '\'' +
				'}';
	}

	@Override
	public boolean equals(Object o) {
		if (this == o) return true;

		if (!(o instanceof APIResource)) {
			return false;
		}

		APIResource that = (APIResource) o;
		return verb.equals(that.verb) &&
				path.equals(that.path);
	}

	@Override
	public int hashCode() {
		return Objects.hash(verb, path);
	}
}


