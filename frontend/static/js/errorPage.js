export function errorPage(msg, code) {
    return /*html*/`
        <div>
            <h1>${msg}</h1>
            <h1>${code}</h1>
        </div>
    `
}