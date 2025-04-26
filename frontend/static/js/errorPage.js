export function errorPage(msg, code) {
    return /*html*/`
        <div class="noPost">
            <h1>${msg}</h1>
            <h1>${code}</h1>
        </div>
    `
}