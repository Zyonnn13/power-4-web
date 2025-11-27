document.addEventListener("DOMContentLoaded", () => {
    
   
    const errorBanner = document.querySelector(".error-banner");

    if (errorBanner) {
       
        errorBanner.animate([
            { transform: 'translateX(0)' },
            { transform: 'translateX(-10px)', offset: 0.1 }, 
            { transform: 'translateX(10px)', offset: 0.2 },
            { transform: 'translateX(-10px)', offset: 0.4 },
            { transform: 'translateX(10px)', offset: 0.6 },
            { transform: 'translateX(0)' }
        ], {
            duration: 500, 
            easing: 'ease-in-out'
        });

        
        setTimeout(() => {
            errorBanner.style.opacity = "0";
            setTimeout(() => {
                errorBanner.remove();
            }, 500); 
        }, 3000); 
    }

   
    if (typeof lastMove !== 'undefined' && lastMove && lastMove.length === 2) {
        const r = lastMove[0]; 
        const c = lastMove[1]; 
        
        const rows = document.querySelectorAll(".board-row");
        if (rows[r]) {
            const cells = rows[r].querySelectorAll(".board-cell");
            const targetCell = cells[c];

            if (targetCell) {
                
                targetCell.style.zIndex = "10"; 
                targetCell.style.position = "relative";

               
                targetCell.animate([
                    { transform: 'translateY(-600px)', opacity: 0, offset: 0 }, 
                    { transform: 'translateY(0)', opacity: 1, offset: 0.6 },    
                    { transform: 'translateY(-40px)', offset: 0.75 },           
                    { transform: 'translateY(0)', offset: 1 }                   
                ], {
                    duration: 600,
                    easing: 'linear',
                    fill: 'forwards'
                });
            }
        }
    }


    if (typeof gameStatus !== 'undefined' && gameStatus === "win") {
        
       
        if (typeof winningCells !== 'undefined' && winningCells && winningCells.length > 0) {
            const rows = document.querySelectorAll(".board-row");
            
            winningCells.forEach(coord => {
                const r = coord[0]; 
                const c = coord[1]; 
                
                const targetRow = rows[r];
                if (targetRow) {
                    const cells = targetRow.querySelectorAll(".board-cell");
                    const targetCell = cells[c];
                    if (targetCell) {
                        
                        targetCell.classList.add("winner-token"); 
                    }
                }
            });
        }

       
        const statusBadge = document.querySelector(".status-badge");
        let timeLeft = 5;

        
        if(statusBadge) statusBadge.textContent = `Victoire ! Redirection dans ${timeLeft}s...`;

        const countdown = setInterval(() => {
            timeLeft--;
            if(statusBadge) statusBadge.textContent = `Victoire ! Redirection dans ${timeLeft}s...`;

            if (timeLeft <= 0) {
                clearInterval(countdown);
                window.location.href = "/game/end"; 
            }
        }, 1000);
    }
    
   
    if (typeof gameStatus !== 'undefined' && gameStatus === "draw") {
         setTimeout(() => {
            window.location.href = "/game/end";
        }, 3000);
    }
});