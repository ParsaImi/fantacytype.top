<script>
  import { onMount } from 'svelte';
  
  const targetText = "سلام اقا پلیسه خوبی";
  let userInput = "";
  let cursorPos = 0;
  let textContainer;
  
  // تقسیم متن به کلمات
  function splitIntoWords(text) {
      return text.split(" ");
  }
  
  // یافتن کلمه و موقعیت کاراکتر فعلی
  function findCurrentWordAndPosition(position) {
      let words = splitIntoWords(targetText);
      let charCount = 0;
      
      for (let i = 0; i < words.length; i++) {
          if (position < charCount + words[i].length) {
              return {
                  wordIndex: i,
                  charIndex: position - charCount,
                  word: words[i]
              };
          }
          
          // اضافه کردن طول کلمه و یک فاصله
          charCount += words[i].length;
          if (i < words.length - 1) {
              charCount++; // برای فاصله
          }
          
          // اگر درست روی فاصله است
          if (position === charCount - 1) {
              return {
                  wordIndex: i,
                  charIndex: words[i].length,
                  word: words[i],
                  isSpace: true
              };
          }
      }
      
      // اگر به انتهای متن رسیده باشیم
      return {
          wordIndex: words.length - 1,
          charIndex: words[words.length - 1].length,
          word: words[words.length - 1]
      };
  }
  
  // بررسی اینکه آیا کاراکتر وارد شده صحیح است
  function isCorrectChar(key, position) {
      return targetText[position] === key;
  }
  
  // به‌روزرسانی نمایش متن
  function updateDisplay() {
      const words = splitIntoWords(targetText);
      const currentPos = findCurrentWordAndPosition(cursorPos);
      
      textContainer.innerHTML = "";
      
      let charCountTotal = 0;
      
      // پردازش هر کلمه
      for (let i = 0; i < words.length; i++) {
          const word = words[i];
          const wordSpan = document.createElement("span");
          wordSpan.className = "word";
          
          // تعیین وضعیت کلمه فعلی
          if (i < currentPos.wordIndex) {
              // کلمه کامل تایپ شده
              const typedPart = userInput.substring(charCountTotal, charCountTotal + word.length);
              const isCorrect = typedPart === word;
              wordSpan.className = `word ${isCorrect ? 'correct' : 'incorrect'}`;
              wordSpan.textContent = word;
          } 
          else if (i === currentPos.wordIndex) {
              // کلمه فعلی که در حال تایپ است
              wordSpan.className = "word current";
              
              // قسمت تایپ شده کلمه فعلی
              for (let j = 0; j < word.length; j++) {
                  const charSpan = document.createElement("span");
                  
                  if (j < currentPos.charIndex) {
                      // کاراکترهای تایپ شده
                      const typedIndex = charCountTotal + j;
                      const isCharCorrect = typedIndex < userInput.length && 
                                            userInput[typedIndex] === word[j];
                      
                      charSpan.className = isCharCorrect ? "correct" : "incorrect";
                      charSpan.textContent = word[j];
                  } 
                  else if (j === currentPos.charIndex) {
                      // کاراکتر فعلی
                      charSpan.className = "current-char";
                      charSpan.textContent = word[j];
                      
                      // اضافه کردن مکان‌نما
                      const cursorSpan = document.createElement("span");
                      cursorSpan.className = "cursor";
                      wordSpan.appendChild(cursorSpan);
                  } 
                  else {
                      // کاراکترهای باقیمانده در کلمه فعلی
                      charSpan.className = "pending";
                      charSpan.textContent = word[j];
                  }
                  
                  wordSpan.appendChild(charSpan);
              }
          } 
          else {
              // کلمات آینده
              wordSpan.className = "word pending";
              wordSpan.textContent = word;
          }
          
          textContainer.appendChild(wordSpan);
          
          // اضافه کردن فاصله بین کلمات
          if (i < words.length - 1) {
              const spaceSpan = document.createElement("span");
              
              if (charCountTotal + word.length === cursorPos) {
                  // مکان‌نما روی فاصله
                  spaceSpan.className = "current-char";
                  spaceSpan.textContent = " ";
                  
                  const cursorSpan = document.createElement("span");
                  cursorSpan.className = "cursor";
                  textContainer.appendChild(cursorSpan);
              } else if (charCountTotal + word.length < cursorPos) {
                  // فاصله تایپ شده
                  const spaceIndex = charCountTotal + word.length;
                  const isSpaceCorrect = spaceIndex < userInput.length && 
                                        userInput[spaceIndex] === " ";
                  
                  spaceSpan.className = isSpaceCorrect ? "correct" : "incorrect";
                  spaceSpan.textContent = " ";
              } else {
                  // فاصله آینده
                  spaceSpan.className = "pending";
                  spaceSpan.textContent = " ";
              }
              
              textContainer.appendChild(spaceSpan);
          }
          
          // به‌روزرسانی شمارنده کاراکتر
          charCountTotal += word.length;
          if (i < words.length - 1) {
              charCountTotal++; // برای فاصله
          }
      }
      
      // اگر به انتهای متن رسیده باشیم، مکان‌نما را نمایش دهیم
      if (cursorPos >= targetText.length) {
          const cursorSpan = document.createElement("span");
          cursorSpan.className = "cursor";
          textContainer.appendChild(cursorSpan);
      }
  }
  
  // بررسی آیا کاراکتر اشتباه وارد شده است
  function hasError() {
      // بررسی کل متن وارد شده تا موقعیت فعلی
      for (let i = 0; i < userInput.length; i++) {
          if (userInput[i] !== targetText[i]) {
              return true;
          }
      }
      return false;
  }
  
  // مدیریت ورودی کیبورد
  function handleKeydown(e) {
      e.preventDefault();
      
      if (e.key === "Backspace") {
          if (cursorPos > 0) {
              cursorPos--;
              userInput = userInput.substring(0, cursorPos);
          }
      } else if (e.key.length === 1) {
          // اگر خطا وجود دارد اجازه حرکت به جلو را نمی‌دهیم
          if (hasError()) {
              // فقط اگر کاراکتر جدید صحیح است و در همان موقعیت قبلی، اجازه تایپ می‌دهیم
              if (isCorrectChar(e.key, cursorPos - 1)) {
                  userInput = userInput.substring(0, cursorPos - 1) + e.key + userInput.substring(cursorPos);
              }
              // در غیر این صورت، کاربر باید ابتدا از Backspace استفاده کند
          } else {
              // اگر خطایی نداریم، اجازه تایپ می‌دهیم
              if (cursorPos < targetText.length) {
                  // فقط اگر کاراکتر صحیح است، اجازه تایپ می‌دهیم
                  if (isCorrectChar(e.key, cursorPos)) {
                      if (cursorPos === userInput.length) {
                          userInput += e.key;
                      } else {
                          userInput = userInput.substring(0, cursorPos) + e.key + userInput.substring(cursorPos);
                      }
                      cursorPos++;
                  }
                  // در غیر این صورت، کاراکتر را نادیده می‌گیریم
              }
          }
      }
      
      updateDisplay();
  }
  
  function handleClick() {
      textContainer.focus();
  }
  
  function resetTyping() {
      userInput = "";
      cursorPos = 0;
      updateDisplay();
      textContainer.focus();
  }
  
  onMount(() => {
      resetTyping();
  });
</script>

<style>
  :global(body) {
      font-family: 'IRANSans', 'Vazir', 'Tahoma', sans-serif;
      text-align: center;
      margin-top: 50px;
      direction: rtl;
      font-size: 24px;
  }
  
  .text-container {
      width: 80%;
      margin: 0 auto;
      padding: 15px;
      border: 1px solid #007bff;
      border-radius: 5px;
      text-align: right;
      line-height: 1.5;
      min-height: 50px;
      position: relative;
      cursor: text;
      word-spacing: 5px;
  }
  
  .text-container:focus {
      outline: 2px solid #007bff;
      box-shadow: 0 0 5px rgba(0, 123, 255, 0.5);
  }
  
  /* کلاس‌های رنگ متن */
  :global(.word) {
      display: inline-block;
      margin: 0 2px;
      white-space: nowrap;
  }
  
  :global(.correct) {
      color: green;
  }
  
  :global(.incorrect) {
      color: red;
  }
  
  :global(.pending) {
      color: gray;
  }
  
  :global(.current) {
      position: relative;
  }
  
  :global(.current-char) {
      background-color: #ffe066;
  }
  
  /* مکان‌نما */
  :global(.cursor) {
      display: inline-block;
      width: 2px;
      height: 1.2em;
      background-color: black;
      animation: blink 1s infinite;
      position: absolute;
      top: 50%;
      transform: translateY(-50%);
  }
  
  @keyframes blink {
      0%, 100% { opacity: 1; }
      50% { opacity: 0; }
  }
  
  /* دکمه شروع مجدد */
  .reset-btn {
      margin-top: 20px;
      padding: 10px 25px;
      font-size: 16px;
      background-color: #f8f9fa;
      border: 1px solid #ced4da;
      border-radius: 4px;
      cursor: pointer;
  }
  
  .reset-btn:hover {
      background-color: #e9ecef;
  }
  
  .error-message {
      color: red;
      margin-top: 10px;
      font-size: 16px;
  }
  
  /* برای حل مشکل فونت */
  @import url('https://cdn.jsdelivr.net/gh/rastikerdar/vazir-font@v30.1.0/dist/font-face.css');
</style>

<div lang="fa" dir="rtl">
  <h2>تمرین تایپ فارسی</h2>
  <p>برای شروع، روی متن کلیک کنید و تایپ کنید</p>
  <p class="error-message">در صورت اشتباه، باید از کلید Backspace استفاده کنید.</p>
  <div 
    bind:this={textContainer} 
    class="text-container" 
    tabindex="0" 
    on:keydown={handleKeydown} 
    on:click={handleClick}>
  </div>
  <button class="reset-btn" on:click={resetTyping}>شروع مجدد</button>
</div>
