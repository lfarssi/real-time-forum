export function createSquare() {
    document.body.addEventListener('mousemove', (e) => {
      const div = document.createElement('div');
      div.classList.add('square');
      div.style.left = e.clientX + 'px';
      div.style.top = e.clientY + 'px';
  
      document.body.append(div);
  
      requestAnimationFrame(() => {
        div.style.opacity = '0';
      });
 
        setTimeout(() => div.remove(), 500);
    })
}
  