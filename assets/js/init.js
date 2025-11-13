const form = document.querySelector('form');

form.setAttribute('novalidate', 'novalidate');

form.addEventListener('submit', formValidation);

function formValidation(event) {
    event.preventDefault();


    const player1Input = form.querySelector('input[name="player1"]');
    const player2Input = form.querySelector('input[name="player2"]');
    
    let isValid = true;


    if (!validateField(player1Input, "Veuillez entrer le nom du joueur 1")) {
        isValid = false;
    }


    if (!validateField(player2Input, "Veuillez entrer le nom du joueur 2")) {
        isValid = false;
    }


    if (isValid) {
        form.submit();
    }
}

function validateField(input, errorMessage) {
    const inputBox = input.closest('.input-box');
    

    const oldError = inputBox.querySelector('.error-message');
    if (oldError) {
        oldError.remove();
    }
    
 
    inputBox.classList.remove('error');
    inputBox.classList.remove('valid');

    if (input.value.trim() === '') {

        const errorDiv = document.createElement('div');
        errorDiv.className = 'error-message';
        errorDiv.textContent = errorMessage;
        inputBox.appendChild(errorDiv);
        
    
        inputBox.classList.add('error');
        
        return false;
    }
    

    inputBox.classList.add('valid');
    return true;
}


const inputs = form.querySelectorAll('input[type="text"]');
inputs.forEach(input => {
    input.addEventListener('input', function() {
        const inputBox = this.closest('.input-box');
        inputBox.classList.remove('error');
        const errorMsg = inputBox.querySelector('.error-message');
        if (errorMsg) {
            errorMsg.remove();
        }
    });
});
