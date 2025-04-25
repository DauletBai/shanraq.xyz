const canvas = document.createElement('canvas');
const ctx = canvas.getContext('2d');
document.getElementById('universe-bg').appendChild(canvas);

// Обновление размера canvas
function resizeCanvas() {
  canvas.width = window.innerWidth;
  canvas.height = window.innerHeight;
}
resizeCanvas();

window.addEventListener('resize', resizeCanvas);

// Создание массива звёзд
let stars = Array(2500).fill().map(() => ({
  x: Math.random() * canvas.width - canvas.width / 2,
  y: Math.random() * canvas.height - canvas.height / 2,
  z: Math.random() * canvas.width,
  color: getRandomColor(),
}));

// Функция для случайного выбора цвета звезды
function getRandomColor() {
  const colors = ['#ffffff', '#ffffcc', '#ddffff', '#ddffdd', '#ffdddd', '#ffddff'];
  return colors[Math.floor(Math.random() * colors.length)];
}

// Функция вычисления прозрачности звезды в зависимости от её положения
function calculateAlpha(x, y, canvasWidth, canvasHeight) {
  const edgeProximity = Math.min(
    x / canvasWidth,
    y / canvasHeight,
    (canvasWidth - x) / canvasWidth,
    (canvasHeight - y) / canvasHeight
  );
  return Math.max(0.9, edgeProximity); // Минимальная прозрачность 0.1
}

// Анимация звёзд с эффектом туманности
function animateStars() {
  // Эффект шлейфа: уменьшаем прозрачность предыдущего кадра
  ctx.fillStyle = 'rgba(0, 0, 0, 0.1)'; // Полупрозрачный чёрный
  ctx.fillRect(0, 0, canvas.width, canvas.height);

  stars.forEach(star => {
    const k = 1500.0 / star.z;
    const x = star.x * k + canvas.width / 2;
    const y = star.y * k + canvas.height / 2;
    const size = Math.max(0, (1 - star.z / canvas.width) * 2);


    // Рассчитываем прозрачность звезды по её положению
    const alpha = calculateAlpha(x, y, canvas.width, canvas.height);

    // Рисуем звезду с градиентной прозрачностью
    if (x >= 0 && x <= canvas.width && y >= 0 && y <= canvas.height) {
      ctx.globalAlpha = alpha; // Устанавливаем прозрачность
      ctx.beginPath();
      ctx.arc(x, y, size, 0, Math.PI * 2);
      ctx.fillStyle = star.color;
      ctx.fill();
    }

    // Обновляем положение звезды
    star.z -= 0.05;
    if (star.z <= 0) {
      star.z = canvas.width;
      star.x = Math.random() * canvas.width - canvas.width / 2;
      star.y = Math.random() * canvas.height - canvas.height / 2;
      star.color = getRandomColor(); // Назначаем новый цвет при перезапуске звезды
    }
  });

  ctx.globalAlpha = 1; // Сбрасываем прозрачность после отрисовки кадра
  requestAnimationFrame(animateStars);
}

animateStars();