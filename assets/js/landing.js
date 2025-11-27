document.addEventListener("DOMContentLoaded", () => {
  const board = document.querySelector("[data-preview-board]");
  if (board) {
    const rows = Array.from(board.querySelectorAll(".board__row"));
    if (rows.length) {
      const colsCount = rows[0].children.length;
      const grid = rows.map((row) => Array.from(row.children));
      const totalRows = grid.length;

      const colors = ["anim-red", "anim-yellow"];
      let colorIndex = 0;

      function resetColumn(col) {
        for (let row = 0; row < totalRows; row++) {
          grid[row][col].classList.remove(
            "anim-red",
            "anim-yellow",
            "rising",
            "bright"
          );
        }
      }

      function launchToken() {
        const col = Math.floor(Math.random() * colsCount);
        resetColumn(col);

        const colorClass = colors[colorIndex];
        colorIndex = (colorIndex + 1) % colors.length;
        const path = [];

        for (let r = totalRows - 1; r >= 0; r--) {
          path.push(grid[r][col]);
        }

        path.forEach((cell, idx) => {
          setTimeout(() => {
            cell.classList.add("rising", colorClass);
            if (idx > 0) {
              path[idx - 1].classList.remove("rising");
            }
            if (idx === path.length - 1) {
              cell.classList.add("bright");
              setTimeout(() => {
                cell.classList.remove("rising");
              }, 120);
              setTimeout(() => {
                cell.classList.remove("bright", colorClass);
              }, 2600);
            }
          }, idx * 90);
        });

        setTimeout(launchToken, 2000);
      }

      launchToken();
    }
  }

  const revealEls = document.querySelectorAll(".reveal");
  if (revealEls.length) {
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            entry.target.classList.add("is-visible");
            observer.unobserve(entry.target);
          }
        });
      },
      {
        threshold: 0.2,
      }
    );

    revealEls.forEach((el) => observer.observe(el));
  }
});

