<script>
  import { onMount } from 'svelte';

  // Text to practice typing
  const targetText = "سلام بر همه حال شما چطور است";
  
  let textContainer;
  let userInput = "";
  let cursorPos = 0;

  // Split text into words
  function splitIntoWords(text) {
    return text.split(" ");
  }

  // Find current word and character position
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
      
      // Add word length and space
      charCount += words[i].length;
      if (i < words.length - 1) {
        charCount++; // For space
      }
      
      // If position is exactly on a space
      if (position === charCount - 1) {
        return {
          wordIndex: i,
          charIndex: words[i].length,
          word: words[i],
          isSpace: true
        };
      }
    }
    
    // If we've reached the end of the text
    return {
      wordIndex: words.length - 1,
      charIndex: words[words.length - 1].length,
      word: words[words.length - 1]
    };
  }

  // Compare user input with target text
  function compareText(input, target, position) {
    let result = [];
    
    for (let i = 0; i < position; i++) {
      if (i < input.length) {
        result.push(input[i] === target[i] ? "correct" : "incorrect");
      } else {
        result.push("pending");
      }
    }
    
    return result;
  }

  // Get display elements for rendering
  function getDisplayElements() {
    const words = splitIntoWords(targetText);
    const currentPos = findCurrentWordAndPosition(cursorPos);
    let elements = [];
    let charCountTotal = 0;
    
    // Process each word
    for (let i = 0; i < words.length; i++) {
      const word = words[i];
      let wordClass = "word";
      let wordContent = [];
      
      // Determine current word state
      if (i < currentPos.wordIndex) {
        // Fully typed word
        const typedPart = userInput.substring(charCountTotal, charCountTotal + word.length);
        const isCorrect = typedPart === word;
        wordClass += ` ${isCorrect ? 'correct' : 'incorrect'}`;
        wordContent.push({ text: word, class: "" });
      } 
      else if (i === currentPos.wordIndex) {
        // Current word being typed
        wordClass += " current";
        
        // Process characters of current word
        for (let j = 0; j < word.length; j++) {
          let charClass = "";
          
          if (j < currentPos.charIndex) {
            // Typed characters
            const typedIndex = charCountTotal + j;
            const isCharCorrect = typedIndex < userInput.length && userInput[typedIndex] === word[j];
            charClass = isCharCorrect ? "correct" : "incorrect";
          } 
          else if (j === currentPos.charIndex) {
            // Current character
            charClass = "current-char";
          } 
          else {
            // Remaining characters in current word
            charClass = "pending";
          }
          
          wordContent.push({ text: word[j], class: charClass });
        }
      } 
      else {
        // Future words
        wordClass += " pending";
        wordContent.push({ text: word, class: "" });
      }
      
      elements.push({
        type: "word",
        class: wordClass,
        content: wordContent
      });
      
      // Add space between words
      if (i < words.length - 1) {
        let spaceClass = "";
        
        if (charCountTotal + word.length < cursorPos) {
          // Typed space
          const spaceIndex = charCountTotal + word.length;
          const isSpaceCorrect = spaceIndex < userInput.length && userInput[spaceIndex] === " ";
          spaceClass = isSpaceCorrect ? "correct" : "incorrect";
        } else {
          // Future space
          spaceClass = "pending";
        }
        
        elements.push({
          type: "space",
          class: spaceClass
        });
      }
      
      // Update character counter
      charCountTotal += word.length;
      if (i < words.length - 1) {
        charCountTotal++; // For space
      }
    }
    
    return {
      elements,
      showCursor: cursorPos <= targetText.length
    };
  }

  // Reset the practice
  function resetPractice() {
    userInput = "";
    cursorPos = 0;
    if (textContainer) {
      textContainer.focus();
    }
  }

  // Handle keyboard input
  function handleKeydown(e) {
    e.preventDefault();
    
    if (e.key === "Backspace") {
      if (cursorPos > 0) {
        cursorPos--;
        userInput = userInput.substring(0, cursorPos);
      }
    } else if (e.key.length === 1) {
      if (cursorPos < targetText.length) {
        // Add character at current position
        if (cursorPos === userInput.length) {
          userInput += e.key;
        } else {
          userInput = userInput.substring(0, cursorPos) + e.key + userInput.substring(cursorPos);
        }
        cursorPos++;
      }
    }
  }

  onMount(() => {
    resetPractice();
  });
</script>

<div class="typing-practice">
  <h2>تمرین تایپ فارسی</h2>
  <p>برای شروع، روی متن کلیک کنید و تایپ کنید</p>
  
  <div 
    bind:this={textContainer}
    class="text-container" 
    tabindex="0" 
    on:keydown={handleKeydown}
    on:click={() => textContainer.focus()}>
    
    {#each getDisplayElements().elements as element}
      {#if element.type === "word"}
        <span class={element.class}>
          {#if element.content.length === 1}
            {element.content[0].text}
          {:else}
            {#each element.content as char}
              <span class={char.class}>{char.text}</span>
            {/each}
          {/if}
        </span>
      {:else if element.type === "space"}
        <span class={element.class}> </span>
      {/if}
    {/each}
    
    {#if getDisplayElements().showCursor}
      <span class="cursor"></span>
    {/if}
  </div>
  
  <button on:click={resetPractice} class="reset-btn">شروع مجدد</button>
</div>

<style>
  .typing-practice {
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
  
  /* Word class styles */
  .word {
    display: inline-block;
    margin: 0 2px;
    white-space: nowrap;
  }
  
  .correct {
    color: green;
  }
  
  .incorrect {
    color: red;
  }
  
  .pending {
    color: gray;
  }
  
  .current {
    position: relative;
  }
  
  .current-char {
    background-color: #ffe066;
  }
  
  /* Cursor style */
  .cursor {
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
  
  /* Reset button */
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
  
  /* Import Vazir font */
  @import url('https://cdn.jsdelivr.net/gh/rastikerdar/vazir-font@v30.1.0/dist/font-face.css');
</style>
