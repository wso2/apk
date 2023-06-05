const vscode = require('vscode');
import * as path from 'path';
const { createConnection, TextDocuments, ProposedFeatures } = require('vscode-languageserver');
const { createServer } = require('yaml-language-server');
// const { Diagnostic, DiagnosticSeverity, Range } = require('vscode-languageserver');

function activate(context) {

console.log('=====================================');

  // Register a document open handler
  context.subscriptions.push(vscode.workspace.onDidOpenTextDocument(handleTextDocumentOpen));

  // Register a document change handler
  // context.subscriptions.push(vscode.workspace.onDidChangeTextDocument(handleTextDocumentChange));

  // Handle opened text documents
  function handleTextDocumentOpen(document) {
    if (document.fileName.endsWith('apk-config.yaml')) {
      // Create a connection to the client
      const connection = createConnection(ProposedFeatures.all);

      // Create a document manager for the opened documents
      const documents = new TextDocuments();

      // Configure the YAML language server
      const server = createServer(connection);

      // Configure the server with the schema file
      const schemaPath = vscode.Uri.file(path.join(context.extensionPath, 'schemas', 'api-config-schema.json')).toString();

      server.configure({
        schemas: [
          {
            fileMatch: ['*.yaml'],
            uri: schemaPath,
          },
        ],
      });

       // Start the server
      server.start(connection);

      // Wire up the server with the document manager
      documents.listen(connection);
      connection.onInitialize(() => {
        return {
          capabilities: {
            textDocumentSync: documents.syncKind,
            completionProvider: {
              resolveProvider: true,
            },
          },
        };
      });
    }
  }

  // function handleTextDocumentChange(change) {
  //   const { document } = change;
  //   if (document.fileName.endsWith('apk-config.yaml')) {
  //     // Validate the document against the schema
  //     validateDocument(document);
  //   }
  // }
}

function deactivate() {}

module.exports = {
  activate,
  deactivate,
};