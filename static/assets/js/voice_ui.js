// static/assets/js/voice_ui.js

document.addEventListener('DOMContentLoaded', () => {
    const aiVisualizer = document.getElementById('ai-visualizer');
    const userVisualizer = document.getElementById('user-visualizer');
    const startButton = document.getElementById('startButton'); // Находим новую кнопку

    // --- Проверка API ... (без изменений) ---
    if (!('speechSynthesis' in window)) { /* ... */ return; }
    const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition;
    if (!SpeechRecognition) { /* ... */ return; }

    const synthesis = window.speechSynthesis;
    const recognition = new SpeechRecognition();
    // ... (настройки recognition без изменений) ...
    recognition.lang = 'ru-RU';
    recognition.interimResults = false;
    recognition.maxAlternatives = 1;


    let voices = [];
    let russianVoice = null;
    let isConversationStarted = false; // Флаг, что диалог начат

    // --- Load Voices (без изменений, но убираем автостарт приветствия) ---
    function populateVoiceList() {
        voices = synthesis.getVoices();
        console.log("Available voices:", voices);
        russianVoice = voices.find(voice => voice.lang === 'ru-RU');
        if (!russianVoice) {
             russianVoice = voices.find(voice => voice.default) || voices[0];
             console.warn("Russian voice not found, using default/first voice:", russianVoice);
        } else {
             console.log("Found Russian voice:", russianVoice);
        }
        // --- НЕ ЗАПУСКАЕМ ПРИВЕТСТВИЕ АВТОМАТИЧЕСКИ ---
        // setTimeout(() => { speak(greeting); }, 500);
    }
    if (synthesis.getVoices().length !== 0) {
        populateVoiceList();
    } else {
        synthesis.onvoiceschanged = populateVoiceList;
        console.log("Waiting for voices to load...");
    }


    // --- Speak Function ---
    function speak(text, callback) { // Добавляем callback
        if (synthesis.speaking) { return; }
        if (text === '') { return; }
        if (!russianVoice) { console.error("No voice selected/available."); return; }

        const utterThis = new SpeechSynthesisUtterance(text);
        // ... (настройки utterThis: voice, lang, pitch, rate) ...
        utterThis.voice = russianVoice;
        utterThis.lang = russianVoice.lang;
        utterThis.pitch = 1;
        utterThis.rate = 1;

        utterThis.onstart = () => {
            console.log('AI Speaking started...');
            aiVisualizer?.classList.add('active');
            userVisualizer?.classList.remove('active');
        };
        utterThis.onend = () => {
            console.log('AI Speaking ended.');
            aiVisualizer?.classList.remove('active');
            if (callback) {
                callback(); // Вызываем callback после окончания речи
            }
        };
        utterThis.onerror = (event) => {
            console.error('SpeechSynthesisUtterance.onerror', event);
            console.error(`Error details: ${event.error}`);
            aiVisualizer?.classList.remove('active');
        };
        synthesis.speak(utterThis);
    }

     // --- Listen Function ---
    function listen() {
        if (!isConversationStarted) return; // Не слушаем, если диалог не начат кнопкой
        console.log("Listening for user...");
        try {
             recognition.start();
        } catch (e) { console.error("recognition.start() error", e); }
    }

    // --- Recognition Event Handlers (без изменений) ---
    recognition.onstart = () => { /* ... */ };
    recognition.onspeechend = () => { /* ... */ };
    recognition.onresult = (event) => { /* ... */ };
    recognition.onerror = (event) => { /* ... */ };
    recognition.onend = () => { /* ... */ };


    // --- Обработчик кнопки СТАРТ ---
    startButton?.addEventListener('click', () => {
        if (isConversationStarted) return; // Не запускать повторно, если уже идет

        isConversationStarted = true; // Начинаем диалог
        console.log("Start button clicked, initiating conversation...");

        // Запрашиваем разрешение на микрофон (неявно через recognition.start())
        // И произносим приветствие. После приветствия - начинаем слушать.
        const greeting = "Здравствуйте! Я ваш семейный AI доктор. Чем могу помочь?";
        speak(greeting, () => {
             // Эта функция будет вызвана, когда приветствие закончится
             console.log("Greeting finished, starting to listen...");
             listen(); // <--- Запускаем распознавание и запрос разрешений
        });

        // Можно добавить изменение вида кнопки после старта, если нужно
        startButton.disabled = true; // Например, деактивировать кнопку
        startButton.innerHTML = '<i class="bi bi-mic"></i> Слушаю...';
    });

});