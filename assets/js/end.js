// Toggle visibility of action buttons around the credits animation
window.addEventListener('DOMContentLoaded', function () {
  const credits = document.querySelector('.credits');
  const actions = document.querySelector('.end-actions');
  if (!credits || !actions) return;

  // Hide actions while credits animation runs
  actions.classList.add('hidden');

  // When animation ends, show the actions again
  credits.addEventListener('animationend', function () {
    actions.classList.remove('hidden');
  });

  // Also show actions immediately if animation is not supported
  setTimeout(() => {
    const style = window.getComputedStyle(credits);
    if (!style.animationName || style.animationName === 'none') {
      actions.classList.remove('hidden');
    }
  }, 100);
});