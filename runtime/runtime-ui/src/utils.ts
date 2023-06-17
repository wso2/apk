/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Red Hat, Inc. All rights reserved.
 *  Copyright (c) Adam Voss. All rights reserved.
 *  Copyright (c) Microsoft Corporation. All rights reserved.
 *  Licensed under the MIT License. See License.txt in the project root for license information.
 *--------------------------------------------------------------------------------------------*/
'use strict';

import { workspace, ExtensionContext, window, commands, Uri } from 'vscode';
import {
  CommonLanguageClient,
  LanguageClientOptions,
  NotificationType,
  RevealOutputChannelOn,
} from 'vscode-languageclient';
import * as path from 'path';
import * as fs from 'fs';
import { SchemaExtensionAPI } from './schema-extension-api';
import { IJSONSchemaCache, JSONSchemaDocumentContentProvider } from './json-schema-content-provider';

export interface ISchemaAssociations {
  [pattern: string]: string[];
}

export interface ISchemaAssociation {
  fileMatch: string[];
  uri: string;
}

// eslint-disable-next-line @typescript-eslint/no-namespace
namespace SchemaAssociationNotification {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  export const type: NotificationType<ISchemaAssociations | ISchemaAssociation[]> = new NotificationType(
    'json/schemaAssociations'
  );
}

// eslint-disable-next-line @typescript-eslint/no-namespace
export namespace SchemaSelectionRequests {
  export const type: NotificationType<void> = new NotificationType('yaml/supportSchemaSelection');
  export const schemaStoreInitialized: NotificationType<void> = new NotificationType('yaml/schema/store/initialized');
}

let client: CommonLanguageClient;
const lsName = 'YAML Support';

export type LanguageClientConstructor = (
  name: string,
  description: string,
  clientOptions: LanguageClientOptions
) => CommonLanguageClient;

export interface RuntimeEnvironment {
  readonly telemetry: TelemetryService;
  readonly schemaCache: IJSONSchemaCache;
}

export interface TelemetryService {
  send(arg: { name: string; properties?: unknown }): Promise<void>;
  sendStartupEvent(): Promise<void>;
}

export function startClient(
  context: ExtensionContext,
  newLanguageClient: LanguageClientConstructor,
  runtime: RuntimeEnvironment
): SchemaExtensionAPI {
  // Options to control the language client
  const clientOptions: LanguageClientOptions = {
    // Register the server for on disk and newly created YAML documents
    documentSelector: [{ language: 'yaml' }, { pattern: '*.y(a)ml' }],
    synchronize: {
      // Notify the server about file changes to YAML and JSON files contained in the workspace
      fileEvents: [workspace.createFileSystemWatcher('**/*.?(e)y?(a)ml'), workspace.createFileSystemWatcher('**/*.json')],
    },
    revealOutputChannelOn: RevealOutputChannelOn.Never,
  };

  // Create the language client and start it
  client = newLanguageClient('yaml', lsName, clientOptions);

  const disposable = client.start();

  const schemaExtensionAPI = new SchemaExtensionAPI(client);


  /////////////////////////// template selection code ////////////////////////////

  const extensionRoot = context.extensionUri.fsPath;
  const templatesFolderPath = Uri.file(path.join(extensionRoot, 'templates'));

  // Register a command to handle the "choose a template" action
  // Read the template files from the templates folder
  const templateFiles = fs.readdirSync(templatesFolderPath.fsPath);

  // Process the template files
  const templates: string[] = [];
  templateFiles.forEach((file) => {
    templates.push(file);
  });

  context.subscriptions.push(
    commands.registerCommand('extension.chooseTemplate', async () => {
      // Show a quick pick menu to let the user choose a template

      const selectedTemplate = await window.showQuickPick(templates);
      if (selectedTemplate) {
        const activeEditor = window.activeTextEditor;
        if (activeEditor && activeEditor.document.fileName.endsWith('.apk-conf')) {
          // Insert the selected template into the currently open 'api-config.yaml' file
          activeEditor.edit((editBuilder) => {
            const templatePath = Uri.file(path.join(templatesFolderPath.fsPath, selectedTemplate));
            const templateContent = fs.readFileSync(templatePath.fsPath, 'utf-8');
            editBuilder.insert(activeEditor.selection.start, templateContent);
          });
        } else {
          window.showErrorMessage('Please open an "apk-config.apk-conf" file.');
        }
      }
    })
  );

  // Show the "choose a template" message when the user creates a new file named "apk-config.yaml"
  context.subscriptions.push(
    workspace.onDidCreateFiles((e) => {
      e.files.forEach((file) => {
        if (file.fsPath.endsWith('.apk-conf')) {
          window.showInformationMessage('Choose a template', 'Select Template').then((choice) => {
            if (choice === 'Select Template') {
              commands.executeCommand('extension.chooseTemplate');
            }
          });
        }
      });
    })
  );
  ////////////////////////// end template selection code ///////////////////////////

  // Push the disposable to the context's subscriptions so that the
  // client can be deactivated on extension deactivation
  context.subscriptions.push(disposable);
  context.subscriptions.push(
    workspace.registerTextDocumentContentProvider(
      'json-schema',
      new JSONSchemaDocumentContentProvider(runtime.schemaCache, schemaExtensionAPI)
    )
  );

  context.subscriptions.push(
    client.onTelemetry((e) => {
      runtime.telemetry.send(e);
    })
  );

  // findConflicts();
  client
    .onReady()
    .then(() => {
      // Send a notification to the server with any YAML schema associations in all extensions
      const associations: ISchemaAssociation[] = [];
      const extensionRootUri = context.extensionUri;
      const schemaUri = Uri.joinPath(extensionRootUri, 'schema', 'apk-schema.json');

      associations.push({
        fileMatch: ['*.apk-conf'],
        uri: schemaUri.toString(),
      });

      client.sendNotification(SchemaAssociationNotification.type, associations);
    })
    .catch((err) => {
      window.showErrorMessage(err);
    });

  return schemaExtensionAPI;
}

export function logToExtensionOutputChannel(message: string): void {
  client.outputChannel.appendLine(message);
}
