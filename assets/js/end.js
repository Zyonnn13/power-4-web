
window.addEventListener('DOMContentLoaded', function () {
  const credits = document.querySelector('.credits');
  const actions = document.querySelector('.end-actions');
  const music = document.getElementById('bgMusic');
  
  if (!credits || !actions) return;

  // Jouer la musique dès le chargement
  if (music) {
    music.volume = 0.5; // Volume à 50%
    music.play().catch(error => {
      console.log("Autoplay bloqué par le navigateur:", error);
      // Si l'autoplay est bloqué, jouer au premier clic
      document.body.addEventListener('click', function playOnClick() {
        music.play();
        document.body.removeEventListener('click', playOnClick);
      }, { once: true });
    });
  }

  actions.classList.add('hidden');

  credits.addEventListener('animationend', function () {
    actions.classList.remove('hidden');
  });

  setTimeout(() => {
    const style = window.getComputedStyle(credits);
    if (!style.animationName || style.animationName === 'none') {
      actions.classList.remove('hidden');
    }
  }, 100);
});