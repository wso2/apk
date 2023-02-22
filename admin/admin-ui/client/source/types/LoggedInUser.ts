
export default interface User {
    name: string,
    _scopes: Array<string>,
    _remember: boolean,
    _environmentName: string,
    rememberMe: boolean
}