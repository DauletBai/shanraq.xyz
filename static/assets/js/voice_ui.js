// static/assets/js/voice_ui.js
document.addEventListener('DOMContentLoaded', () => {
    const aiVisualizer = document.getElementById('ai-visualizer');
    const userVisualizer = document.getElementById('user-visualizer');
    const startButton = document.getElementById('startButton');

    console.log("Voice UI Script Loaded.");

    // --- Проверка API ---
    if (!('speechSynthesis' in window)) { console.error("Browser does not support Speech Synthesis."); return; }
    const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition;
    if (!SpeechRecognition) { console.error("Browser does not support Speech Recognition."); return; }

    const synthesis = window.speechSynthesis;
    const recognition = new SpeechRecognition();
    // --- Настройки Recognition ---
    recognition.lang = 'ru-RU';
    recognition.interimResults = false; // Оставляем false
    recognition.maxAlternatives = 1;
    recognition.continuous = false;     // Оставляем false (одна фраза за раз)

    console.log("SpeechRecognition object created.");

    let voices = [];
    let russianVoice = null;
    let isConversationStarted = false;

    // --- Load Voices ---
    function populateVoiceList() {
        voices = synthesis.getVoices();
        console.log("Available voices:", voices.length);
        russianVoice = voices.find(voice => voice.lang === 'ru-RU');
        if (!russianVoice) {
             russianVoice = voices.find(voice => voice.default) || voices[0];
             console.warn("Russian voice not found, using:", russianVoice?.name);
        } else {
             console.log("Found Russian voice:", russianVoice.name);
        }
    }
    if (synthesis.getVoices().length !== 0) { populateVoiceList(); }
    else { synthesis.onvoiceschanged = populateVoiceList; console.log("Waiting for voices to load..."); }

    // --- Speak Function ---
    function speak(text, callback) {
         if (synthesis.speaking) { console.warn("Synthesis busy."); return; }
         if (!russianVoice) { console.error("No voice selected for speak."); return; }
         // console.log(`Attempting to speak: "${text.substring(0, 30)}..."`); // Убрал для чистоты
         const utterThis = new SpeechSynthesisUtterance(text);
         utterThis.voice = russianVoice;
         utterThis.lang = russianVoice.lang;
         utterThis.pitch = 1;
         utterThis.rate = 1;
         utterThis.onstart = () => { console.log('AI Speaking started...'); aiVisualizer?.classList.add('active'); userVisualizer?.classList.remove('active'); };
         utterThis.onend = () => { console.log('AI Speaking ended.'); aiVisualizer?.classList.remove('active'); if (callback) { callback(); } };
         utterThis.onerror = (event) => { console.error('SpeechSynthesisUtterance.onerror:', event.error); aiVisualizer?.classList.remove('active'); };
         synthesis.speak(utterThis);
    }

    // --- Listen Function ---
    function listen() {
        if (!isConversationStarted) return;
        console.log("Attempting to start recognition...");
        try {
             recognition.start();
        } catch (e) { console.error("recognition.start() error:", e); }
    }

    // --- Recognition Event Handlers ---
    recognition.onstart = () => {
        console.log("Event: recognition.onstart fired");
        userVisualizer?.classList.add('active');
        console.log("User visualizer activated.");
    };

    recognition.onspeechend = () => {
        console.log("Event: recognition.onspeechend fired");
        // Ничего не делаем
    };

    recognition.onresult = (event) => {
        console.log("Event: recognition.onresult fired");
        // console.log("Result event:", event); // Можно закомментировать для чистоты лога
        if (event.results.length > 0 && event.results[event.results.length - 1].isFinal) {
             let transcript = event.results[event.results.length - 1][0].transcript.trim();
             console.log(`Final transcript received: "${transcript}"`);
             // !!! ЗДЕСЬ НУЖНО БУДЕТ ОТПРАВЛЯТЬ transcript НА БЭКЕНД !!!
             // Например, через WebSocket: websocket.send(transcript);

             // Пока просто выводим в консоль и ничего не делаем дальше
             // Остановка произойдет автоматически и вызовется onend
        }
    };

    recognition.onerror = (event) => {
        console.log("Event: recognition.onerror fired");
        console.error('Recognition error:', event.error);
        userVisualizer?.classList.remove('active'); // Убираем анимацию при ошибке
    };

    recognition.onend = () => {
        console.log("Event: recognition.onend fired");
        userVisualizer?.classList.remove('active'); // Убираем анимацию при завершении
         if (isConversationStarted) {
             startButton.disabled = false;
             startButton.innerHTML = '<i class="bi bi-mic-fill me-2"></i> Начать консультацию';
             isConversationStarted = false; // Сбрасываем флаг, чтобы можно было начать снова
             console.log("Recognition service ended. Ready for next interaction.");
         }
    };

    // --- Start Button Handler ---
    startButton?.addEventListener('click', () => {
         if (!russianVoice) { alert("Голосовой движок не готов."); return; }
         if (isConversationStarted) { console.warn("Conversation already started."); return; } // Предотвращаем двойной старт
         isConversationStarted = true;
         console.log("Start button clicked.");
         const greeting = "Здравствуйте! Я ваш семейный AI доктор. Чем могу помочь?";
         speak(greeting, () => {
             console.log("Greeting finished callback: starting listen().");
             listen();
         });
         startButton.disabled = true;
         startButton.innerHTML = '<i class="bi bi-mic"></i> Слушаю...';
    });

});