import * as vscode from 'vscode';
import * as fs from 'fs';
import * as path from 'path';

export async function activate(context: vscode.ExtensionContext) {

	// Register the "*.apk-conf" file association to the "yaml" language
	// registerFileAssociation();

	// Check if "YAML Language Support by Red Hat" extension is installed
	const yamlExtension = vscode.extensions.getExtension('redhat.vscode-yaml');
	if (!yamlExtension) {
		vscode.window.showErrorMessage(
			'The "YAML Language Support by Red Hat" extension is required for the APK Configuration extension to work properly. Please install it and reload the window.'
		);
		return;
	}
	const yamlExtensionAPI = await yamlExtension.activate();
	const SCHEMA = "apkschema";

	// Read the schema file content
	const schemaFilePath = path.join(context.extensionPath, 'schema', 'apk-schema.json');

	const schemaContent = fs.readFileSync(schemaFilePath, 'utf8');
	const schemaContentJSON = JSON.parse(schemaContent);

	const schemaJSON = JSON.stringify(schemaContentJSON);

	/**
	 * 
	 * @param resource  The URI of the resource
	 * @returns  The URI of the schema file
	 * 
	 * This function is called when the YAML Language Support extension needs to know the URI of the schema file.
	 * The schema file is stored in the extension's "schema" folder.
	 * The schema file is named "apk-schema.json".
	 */
	function onRequestSchemaURI(resource: string): string | undefined {
		if (resource.endsWith('.apk-conf')) {
			return `${SCHEMA}://schema/apk-conf`;
		}
		return undefined;
	}
	/**
	 * 
	 * @param schemaUri  The URI of the schema
	 * @returns  The content of the schema file
	 *
	 */
	function onRequestSchemaContent(schemaUri: string): string | undefined {
		const parsedUri = vscode.Uri.parse(schemaUri);
		if (parsedUri.scheme !== SCHEMA) {
			return undefined;
		}
		if (!parsedUri.path || !parsedUri.path.startsWith('/')) {
			return undefined;
		}

		return schemaJSON;
	}

	// Register the schema provider
	yamlExtensionAPI.registerContributor(SCHEMA, onRequestSchemaURI, onRequestSchemaContent);


	/////////////////////////// template selection code ////////////////////////////

	const extensionRoot = context.extensionUri.fsPath;
	const templatesFolderPath = vscode.Uri.file(path.join(extensionRoot, "templates"));

	/**
	 * Register a command to handle the "choose a template" action
	 * Read the template files from the templates folder
	 * Process the template files
	 */
	// vscode.workspace.onDidOpenTextDocument((document) => {

	// Register a command to handle the "choose a template" action
	// Read the template files from the templates folder
	const templateFiles = fs.readdirSync(templatesFolderPath.fsPath);

	// Process the template files
	const templates: string[] = [];
	templateFiles.forEach((file) => {
		templates.push(file);
	});

	context.subscriptions.push(
		vscode.commands.registerCommand("extension.chooseTemplateApk", async () => {
			// Show a quick pick menu to let the user choose a template
			const selectedTemplate = await vscode.window.showQuickPick(templates);
			if (selectedTemplate) {
				const activeEditor = vscode.window.activeTextEditor;
				if (
					activeEditor &&
					activeEditor.document.fileName.endsWith(".apk-conf")
				) {
					// Insert the selected template into the currently open '*.apk-conf' file
					activeEditor.edit((editBuilder) => {
						const templatePath = vscode.Uri.file(
							path.join(templatesFolderPath.fsPath, selectedTemplate)
						);
						const templateContent = fs.readFileSync(
							templatePath.fsPath,
							"utf-8"
						);
						editBuilder.insert(activeEditor.selection.start, templateContent);
					});
				} else {
					vscode.window.showErrorMessage('Please open an "apk-config.apk-conf" file.');
				}
			}

		})
	);

	// Show the "choose a template" message when the user creates a new file named "apk-config.yaml"
	// in the workspace
	context.subscriptions.push(
		vscode.workspace.onDidCreateFiles((e) => {
			e.files.forEach((file) => {
				if (file.fsPath.endsWith(".apk-conf")) {
					vscode.window
						.showInformationMessage("Choose a template", "Select Template")
						.then((choice) => {
							if (choice === "Select Template") {
								vscode.commands.executeCommand("extension.chooseTemplateApk");
							}
						});
				}
			});
		})
	);
	// });
	////////////////////////// end template selection code ///////////////////////////
}


