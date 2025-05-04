export function errorPage(msg) {
    return /*html*/`
        <div class="noPost">
            <h1>${msg}</h1>
        </div>
    `
}
function throttle(fn, delay=2000) {
    let inThrottle = false;
    return function (...args) {
      if (inThrottle) return;
      fn.apply(this, args);
      inThrottle = true;
      setTimeout(() => inThrottle = false, delay);
    };
  }
  
  
  export function popup(msg, status) {
    const popupDiv = document.createElement('div');
    popupDiv.className = 'popup';

    const icon = status === 'success'
        ? '<i class="fa-solid fa-circle-check"></i>'
        : '<i class="fa-solid fa-circle-exclamation"></i>';

    const iconColor = status === 'success' 
        ? '#4caf50' 
        : '#f44336';

    const closeButtonClass = status === 'success' ? 'popup-close-success' : 'popup-close-failed';
    const timerClass = status === 'success' ? 'popup-timer-success' : 'popup-timer-failed';
    
    popupDiv.innerHTML = `
      <span class="popup-icon" style="color: ${iconColor}; font-size: 20px; margin-right: 10px;">${icon}</span>
      <span class="popup-message">${msg}</span>
      <button class="${closeButtonClass} popup-close">✖</button>
      <div class="${timerClass} popup-timer"></div>
    `;

    document.body.appendChild(popupDiv);

    const closeButton = popupDiv.querySelector('.popup-close');
    closeButton.addEventListener('click', () => {
        popupDiv.remove(); 
    });

    setTimeout(() => {
        popupDiv.remove();
    }, 3000);
}


export const popupThrottled = throttle(popup, 3000);
