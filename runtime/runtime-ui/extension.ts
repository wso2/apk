import * as vscode from 'vscode';
import * as path from 'path';
import { TextDocument, CompletionItem, CompletionItemKind } from 'vscode';

function provideCompletionItems(document: TextDocument): CompletionItem[] {
    // Load the schema file
    const schemaPath = path.join(__dirname, 'schemas', 'api-schema.yaml');
    const schemaUri = vscode.Uri.file(schemaPath);

    // Get the JSON representation of the schema
    const schema = vscode.workspace.getConfiguration().get('yaml', {}).get('schemas', []);
    const schemaConfiguration = {
        fileMatch: [document.uri.toString()],
        schema: {
            uri: schemaUri.toString(),
        },
    };
    schema.push(schemaConfiguration);
    vscode.workspace.getConfiguration().update('yaml.schemas', schema, true);

    // Provide completion items based on the schema
    const completionItems: CompletionItem[] = [];

    // TODO: Implement your code completion logic here based on the schema

    return completionItems;
}

export function activate(context: vscode.ExtensionContext) {
    const provider = vscode.languages.registerCompletionItemProvider('yaml', {
        provideCompletionItems,
    });

    context.subscriptions.push(provider);
}
