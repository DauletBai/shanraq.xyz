// Вызываем определение элементов LUME как можно раньше,
// чтобы браузер знал о тегах lume-*, когда будет парсить HTML.
// Этот вызов может быть и здесь, или даже раньше, если подключить lume до earth.js
LUME.defineElements();

// Весь остальной код, работающий с DOM элементами,
// помещаем внутрь обработчика события DOMContentLoaded.
document.addEventListener('DOMContentLoaded', (event) => {
    // Получаем ссылки на элементы *после* того, как DOM гарантированно загружен
    const earth = document.getElementById('earth');
    const clouds = document.getElementById('clouds');
    const moonRotator = document.getElementById('moonRotator');

    // Проверяем, найдены ли элементы (на всякий случай)
    if (!earth || !clouds || !moonRotator) {
        console.error("One or more LUME elements (earth, clouds, moonRotator) not found!");
        return;
    }

    // Теперь назначаем анимацию
    let lastTime = performance.now();
    let dt = 0;

    // Вращение Луны
    moonRotator.rotation = (x, y, z, time) => {
        dt = time - lastTime;
        lastTime = time;
        // Вращаем вокруг оси Z узла moonRotator
        return [x, y, z + dt * 0.001];
    };

    // Вращение Земли
    earth.rotation = (x, y, z, t) => [x, t * 0.001, z];

    // Вращение облаков
    clouds.rotation = (x, y, z, t) => [x, -t * 0.001, z]; // Небольшая рассинхронизация с Землей

    console.log("Earth rotation logic applied."); // Отладочное сообщение
});
