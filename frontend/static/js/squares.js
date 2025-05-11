export function createSquare() {
    document.body.removeEventListener('mousemove', squareMouseHandler);
    document.body.addEventListener('mousemove', squareMouseHandler);
}

export function squareMouseHandler(e) {
    if (e.clientX + 20 < innerWidth && e.clientY + 20 < innerHeight) {
        const div = document.createElement('div');
        div.classList.add('square');
        div.style.left = e.clientX + 'px';
        div.style.top = e.clientY + 'px';

        document.body.append(div);

        requestAnimationFrame(() => {
            div.style.opacity = '0';
        });

        setTimeout(() => div.remove(), 500);
    }
}
