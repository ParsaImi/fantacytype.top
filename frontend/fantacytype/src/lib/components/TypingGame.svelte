<script lang="ts">
  import { onMount, onDestroy, tick, createEventDispatcher } from 'svelte';
  import { goto } from '$app/navigation';
  

  export let gameData: {
    text: string,
    startTime: string,
    IsActive: string,
    language?: string // Add language property
  };
  export let onWordComplete: (word: string) => void;
  
  // Add event dispatcher to communicate with parent
  
  type PlayerResult = { 
      type: string;
      playerid; number;
      rank: number;
      wpm: number;
      durationminutes: number;
  }
  const dispatch = createEventDispatcher();
  let playerserverid: string;
  let textContainer;
  let userInput = "";
  let errorCount = 0;
  let totalTyped = 0;
  let cursorPos = 0;
  let gameText = '';
  let textArray: string[] = [];
  let currentPosition = 0;
  let typedText = '';
  
  let isActive = false;
  let startTime: Date;
  let currentWordStart = 0;
  let wordsTyped = 0;
  let accuracy = 100;
  let correctChars = 0;
  let totalChars = 0;
  let errorState = false; // Flag to track if user is in error state
  let isPersian = false; // Flag for Persian text
  
  // Progress tracking - from server only
  let playerProgress: Record<string, number> = {};
  let currentUserId = ''; // This would be set from localStorage or session
  
  // Player ranks tracking
  let playerRanks: Record<string, PlayerResult> = {};
  
  let countdownTimeLeft = 0;
  let countdownInterval: number;
  let gameStarted = false;
  // Game end state
  let gameEnded = false;
  let gameResults: any = null;
  let finalLeaderboard: Record<number, string> = {};
  let token = '';
  let isAuthenticated = false;
  let username = '';
  
  $: if (gameData) {
    gameText = gameData.text;
    textArray = gameText.split('');
    isActive = "FALSE";
    startTime = new Date(parseInt(gameData.startTime));
    currentPosition = 0;
    typedText = '';
    currentWordStart = 0;
    wordsTyped = 0;
    correctChars = 0;
    totalChars = 0;
    errorState = false;
    gameEnded = false;
    // Set Persian flag if language is specified
    isPersian = gameData.language === 'persian';
    // Reset player ranks
    playerRanks = {};
    finalLeaderboard = {};
    startCountdown();
  }

  
  $: if (totalTyped> 0) {
    accuracy = Math.round(((totalTyped - errorCount) / totalTyped) * 100);
  }

  function startCountdown() {
      // Guard clause to prevent countdown if component is destroyed
      
      console.log("countdown started!!", startTime)
      console.log(gameData.startTime)
      if (countdownInterval) {
          clearInterval(countdownInterval)
      }
      const now = new Date();
      console.log("NOWWWWW")
      console.log(now)
      const timeDiff = startTime.getTime() - now.getTime();

      if (timeDiff <= 0) {
          console.log(`BA KAMAL MEYL  ${timeDiff}`)
          countdownTimeLeft = 0;
          gameEnded = false
          startGame();
          return;
      }

      countdownTimeLeft = Math.ceil(timeDiff / 1000)

      countdownInterval = setInterval(() => {
          
          
          countdownTimeLeft--;

          if (countdownTimeLeft <= 0) {
              clearInterval(countdownInterval)
              gameEnded = false
              startGame()
          }
          console.log("time goes!!")
        }, 1000)
    }
    
  function startGame() { // Guard clause to prevent starting game if component is destroyed
      
      
      gameStarted = true;
      isActive = true

      setTimeout(() => {
          if (textContainer){
                textContainer.focus();
                updateDisplay()
          }
          
          }, 100)
  }
  
  function hasError() {
      // بررسی کل متن وارد شده تا موقعیت فعلی
      for (let i = 0; i < userInput.length; i++) {
          if (userInput[i] !== gameText[i]) {
              return true;
          }
      }
      return false;
  } 
  
 function handleKeydown(e) {
    // Guard clause to prevent execution if component is destroyed or game is not active
    if (!textContainer || !gameText || !isActive || gameEnded) {
        return;
    }
    
    console.log(e.key)
    e.preventDefault();
    
    if (e.key === "Backspace") {
        // Only allow backspace when in error state
        if (errorState) {
            // In error state, just clear the error state and let the user try again
            // Don't move cursor position and don't modify userInput length
            // We'll just truncate userInput to match the cursor position
            userInput = gameText.substring(0, cursorPos);
            errorState = false;
            updateDisplay();
        }
        // Ignore backspace if not in error state
    } else if (e.key.length === 1) {
        if (cursorPos < gameText.length) {
            // Check if key matches the expected character at cursor position
            if (e.key === gameText[cursorPos]) {
                // Correct character - update user input
                if (cursorPos === userInput.length) {
                    userInput += e.key;
                } else {
                    userInput = userInput.substring(0, cursorPos) + e.key + userInput.substring(cursorPos);
                }

                totalTyped++;
                errorState = false; // Reset error state when typing correct character
                
                // CHECK FOR COMPLETED WORD
                if (e.key === ' ' || cursorPos + 1 === gameText.length) {
                    const completedWord = gameText.substring(currentWordStart, cursorPos + 1).trim();
                    console.log(completedWord);
                    if (completedWord) {
                        onWordComplete(completedWord);
                        currentWordStart = cursorPos + 1;
                        wordsTyped++;
                    }
                }
                
                // Increment cursor position for correct character
                cursorPos++;
            } else {
                // Incorrect character - record the error but don't advance cursor
                errorCount++;
                totalTyped++;
                errorState = true;
                // Note: We don't update userInput here, keeping the player at the same position
            }
        }
    }
    
    updateDisplay();
} 

  // Cleanup function to remove event listeners and reset variables
  function cleanup() {
    window.removeEventListener('keydown', handleKeydown);
  }
  
  function handleClick() {
      textContainer.focus();
  }
  
  onMount(() => {
    window.addEventListener('keydown', handleKeydown);
    
    // Get username from localStorage if available
    token = localStorage.getItem('auth_token');
    if (token) {
      try {
        // Basic validation and extract username from token
        const parts = token.split('.');
        if (parts.length === 3) {
          // Decode the payload (middle part of the JWT)
          const payload = JSON.parse(atob(parts[1]));
          if (payload && payload.sub) {
            username = payload.sub;
            currentUserId = username
            isAuthenticated = true;
          }
        }
      } catch (error) {
        console.error('Error decoding token:', error);
      }
    }else{
        username = localStorage.getItem('user_nickname') || '';
        
        currentUserId = username
        if (!username) {
          // Redirect to nickname page if no nickname is set
          goto('/nickname?redirect=/game');
          return;
        }
    }  
    updateDisplay()
    
    if (gameData && gameData.startTime) {
        console.log("startCountdown goes")
        startCountdown();
    } else {
        console.log("NOT THE COUNTER TIMEj")
    }
    
    // Return cleanup function
    return cleanup;
  });
  
  // Use onDestroy to ensure cleanup always happens
  onDestroy(() => {
    cleanup();
  });
  
  // Handle new websocket message for progress updates
  export function updateProgress(userId: string, percentage: number) {
    playerProgress[userId] = percentage;
    playerProgress = {...playerProgress}; // Trigger reactivity
    console.log(userId + " IDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD")
  }
  
  // New function to handle player rank updates
  export function updatePlayerRanks(data: PlayerResult) {
    playerRanks[data.playerid] = data
    playerserverid = data.playerid
    console.log("HERE SHITTTTTTTTTT")
    console.log(currentUserId)
    console.log(data)
  }
  
  // Handle game end signal
  export function endGame(results = null) {
    gameStarted = false
    cursorPos = 0
    clearInterval(countdownInterval)
    typedText = ""
    userInput = ""
    gameEnded = true;
    gameResults = results;
    
    // Process leaderboard data if available
    console.log(`this is the final result ${results} and ${results.leaderboard}`)
    if (results && results.leaderboard) {
      finalLeaderboard = results.leaderboard;
    }
    errorCount = 0;
    totalTyped = 0;

    isActive = false;
  }
  
  function getCharClass(char: string, index: number) {
    if (index >= typedText.length) {
      return 'not-typed';
    }
    
    return typedText[index] === char ? 'correct' : 'incorrect';
  }
  
  // Get rank badge color based on position
  function getRankBadgeColor(rank: number) {
    switch(rank) {
      case 1: return 'gold';
      case 2: return 'silver';
      case 3: return 'bronze';
      default: return 'default';
    }
  }
  
  // Sort players by progress for the leaderboard
  $: sortedPlayers = Object.entries(playerProgress)
    .sort(([, progressA], [, progressB]) => progressB - progressA);
    
  // Get formatted leaderboard entries from final results
  $: leaderboardEntries = Object.entries(finalLeaderboard)
    .map(([rank, username]) => ({
      rank: parseInt(rank),
      username: username,
      isCurrentUser: username === currentUserId,
      wpm: playerRanks[playerserverid]?.wpm || 0,
      duration: playerRanks[playerserverid]?.durationminutes || 0,
    }))
    .sort((a, b) => a.rank - b.rank);

  
    
  // Navigation functions
  function goToMainMenu() {
    goto('/');
  }
  
  function playAgain() {
    // Instead of redirecting, dispatch an event to tell parent component to go back to waiting
    window.location.reload();
  }
  
  // Get player rank from the results
  function getCurrentPlayerRank() {
    // Find the player in the final leaderboard
    for (const [rank, username] of Object.entries(finalLeaderboard)) {
      if (username === currentUserId) {
        return parseInt(rank);
      }
    }
    
    // Fallback to real-time rank or default
    return playerRanks[currentUserId] || 
      (sortedPlayers.findIndex(([playerId]) => playerId === currentUserId) + 1);
  }

  function splitIntoWords(text){
      return text.split(" ")
  }

  function findCurrentWordAndPosition(position) {
      let words = splitIntoWords(gameText);
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

  function isCorrectChar(key, position) {
      return gameText[position] === key;
  }

  function updateDisplay() {
      console.log(cursorPos)
      console.log(playerRanks)
      console.log("updating Display !!!!!")
      const words = splitIntoWords(gameText);
      console.log(words)
      const currentPos = findCurrentWordAndPosition(cursorPos);
      
      textContainer.innerHTML = "";
      
      let charCountTotal = 0;
      
      // پردازش هر کلمه
      for (let i = 0; i < words.length; i++) {
          const word = words[i];
          const wordSpan = document.createElement("span");
          wordSpan.className = "word";
          
          // DETERMINING CURRENT WORD SITUATION
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
                  console.log(errorState)
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
                      charSpan.className = errorState ? "current-char error" : "current-char";
                      charSpan.textContent = word[j];
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
              
              if (charCountTotal + word.length === cursorPos && errorState) {
                  spaceSpan.className = "space error";
                  spaceSpan.textContent = " ";
              }
              else if (charCountTotal + word.length === cursorPos) {
                  // مکان‌نما روی فاصله
                  spaceSpan.className = "current-char";
                  spaceSpan.textContent = " ";
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
  }

  function formatCountdown(seconds: number): string {
      const minutes = Math.floor(seconds / 60)
      const remainingSeconds = seconds % 60;
      return `${minutes.toString().padStart(2, '0')}:${remainingSeconds.toString().padStart(2, '0')}`;
  }
  
</script>

<div class="game-container">
  <h2>Typing Game {isPersian ? '(Persian)' : ''}</h2>
  
  {#if gameEnded}
    <div class="game-end-popup">
      <div class="popup-content">
        <h3>Game Ended</h3>
        
        <div class="results">
          <p>Your final score:</p>
          <ul>
            <li>Words typed: {wordsTyped}</li>
            <li>Accuracy: {accuracy}%</li>
            <li>Final position: {getCurrentPlayerRank()} of {Object.keys(finalLeaderboard).length || sortedPlayers.length}</li>
          </ul>
        </div>
        
        <!-- Final leaderboard display -->
        <div class="final-leaderboard">
          <h4>Final Leaderboard</h4>
          <table>
              <thead>
                <tr>
                  <th>Rank</th>
                  <th>Player</th>
                  <th>WPM</th>
                  <th>Duration</th>
                </tr>
              </thead>
              <tbody>
                {#if leaderboardEntries.length > 0}
                  {#each leaderboardEntries as entry}
                    <tr class={entry.isCurrentUser ? 'current-user' : ''}>
                      <td>
                        <div class="rank-badge rank-{getRankBadgeColor(entry.rank)}">
                          #{entry.rank}
                        </div>
                      </td>
                      <td>{entry.username} {entry.isCurrentUser ? '(You)' : ''}</td>
                      <td class="wpm-cell">{entry.wpm.toFixed(1)}</td>
                      <td class="duration-cell">{entry.duration.toFixed(1)}m</td>
                    </tr>
                  {/each}
                {:else}
                  {#each sortedPlayers as [playerId, progress], index}
                    <tr class={playerId === currentUserId ? 'current-user' : ''}>
                      <td>
                        <div class="rank-badge rank-{getRankBadgeColor(index + 1)}">
                          #{index + 1}
                        </div>
                      </td>
                      <td>{playerId} {playerId === currentUserId ? '(You)' : ''}</td>
                      <td class="wpm-cell">-</td>
                      <td class="duration-cell">-</td>
                    </tr>
                  {/each}
                {/if}
              </tbody>
          </table> 
        </div>
        
        <div class="popup-buttons">
          <button class="btn-main" on:click={goToMainMenu}>Main Menu</button>
          <button class="btn-play" on:click={playAgain}>Play Again</button>
        </div>
      </div>
    </div>
  {/if}
  
  <!-- Countdown Timer -->
  {#if !gameStarted && !gameEnded && countdownTimeLeft > 0}
    <div class="countdown-container">
      <h3>Game Starting In</h3>
      <div class="countdown-timer">{formatCountdown(countdownTimeLeft)}</div>
      <p>Get ready to type!</p>
    </div>
  {:else if isActive && !gameEnded}
    <div class="stats">
      <div class="stat">
        <span class="label">Words:</span>
        <span class="value">{wordsTyped}</span>
      </div>
      <div class="stat">
        <span class="label">Accuracy:</span>
        <span class="value">{accuracy}%</span>
      </div>
      <div class="stat">
        <span class="label">Progress:</span>
        <span class="value">{playerProgress[currentUserId] || 0}%</span>
      </div>
      {#if playerRanks[currentUserId]}
        <div class="stat">
          <span class="label">Rank:</span>
          <span class="value">{playerRanks[currentUserId]}</span>
        </div>
      {/if}
    </div>

    <!-- Error state notification -->
    {#if errorState}
      <div class="error-message">
        Type the correct character to continue
      </div>
    {/if}
    
    <!-- Progress race track - based entirely on server data -->
    <div class="progress-race">
      <h3>Race Progress</h3>
      <div class="progress-container">
        {#each sortedPlayers as [playerId, progress]}
          <div class="player-progress">
            <div class="player-name">{playerId === currentUserId ? `${currentUserId}  (You)` : playerId}</div>
            <div class="progress-bar-container">
              <div class="progress-bar" style="width: {progress}%"></div>
              <div class="progress-value">{progress}%</div>
            </div>
            {#if playerRanks[playerId]}
              <div class="rank-badge rank-{getRankBadgeColor(playerRanks[playerId])}">
                #{playerRanks[playerId]}
              </div>
            {/if}
          </div>
        {/each}
      </div>
    </div>
    
    <div bind:this={textContainer} class="text-container" class:error-state={errorState} class:rtl={isPersian} tabindex="0">
    </div>
  {:else if !gameEnded}
    <div class="waiting">
      <p>Waiting for the game to start...</p>
    </div>
  {/if}
</div>

<style>
  /* Dark Mode Styling - Modern UI for Typing Game */
  
  :global(body) {
    font-family: 'IRANSans', 'Vazir', 'Tahoma', 'Segoe UI', -apple-system, BlinkMacSystemFont, sans-serif;
    background-color: #121212;
    color: #e0e0e0;
    text-align: center;
    margin: 0;
    padding: 0;
    min-height: 100vh;
    transition: background-color 0.3s ease;
  }
  
  .game-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 2rem;
    max-width: 1200px;
    margin: 0 auto;
  }
  
  h2 {
    font-size: 2.5rem;
    margin-bottom: 2rem;
    text-align: center;
    color: #bb86fc; /* Purple accent */
    font-weight: 600;
    letter-spacing: -0.5px;
  }
  
  h3 {
    font-size: 1.5rem;
    margin-bottom: 1rem;
    text-align: center;
    color: #bb86fc;
    font-weight: 500;
  }
  
  h4 {
    font-size: 1.2rem;
    margin-bottom: 0.75rem;
    text-align: center;
    color: #e0e0e0;
  }
  
  /* Stat display */
  .stats {
    display: flex;
    justify-content: space-around;
    width: 100%;
    max-width: 800px;
    margin-bottom: 2rem;
    gap: 1rem;
  }
  
  .stat {
    padding: 1rem;
    background-color: #1e1e1e;
    border-radius: 8px;
    display: flex;
    flex-direction: column;
    align-items: center;
    min-width: 100px;
    flex: 1;
    border: 1px solid #333;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    transition: transform 0.2s ease, box-shadow 0.2s ease;
  }
  
  .stat:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 10px rgba(0, 0, 0, 0.15);
  }
  
  .label {
    font-size: 0.9rem;
    color: #aaa;
    margin-bottom: 0.3rem;
  }
  
  .value {
    font-size: 1.4rem;
    font-weight: bold;
    color: #fff;
  }
  
  /* Error message */
  .error-message {
    background-color: rgba(255, 82, 82, 0.15);
    color: #ff5252;
    padding: 1rem 1.5rem;
    border-radius: 8px;
    margin-bottom: 1.5rem;
    font-weight: 500;
    border-left: 4px solid #ff5252;
    width: 100%;
    max-width: 800px;
    text-align: center;
    backdrop-filter: blur(4px);
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    animation: pulse 1.5s infinite;
  }
  
  @keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.7; }
  }
  
  /* Progress race */
  .progress-race {
    width: 100%;
    max-width: 800px;
    margin-bottom: 2rem;
    padding: 1.5rem;
    background-color: #1e1e1e;
    border-radius: 12px;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
  }
  
  .progress-container {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin-top: 1rem;
  }
  
  .player-progress {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 0.5rem;
    border-radius: 8px;
    background-color: #252525;
    transition: background-color 0.2s ease;
  }
  
  .player-progress:hover {
    background-color: #2c2c2c;
  }
  
  .player-name {
    flex: 0 0 120px;
    font-weight: 500;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    color: #e0e0e0;
    text-align: left;
  }
  
  .progress-bar-container {
    flex: 1;
    height: 24px;
    background-color: #333;
    border-radius: 12px;
    overflow: hidden;
    position: relative;
    box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.3);
  }
  
  .progress-bar {
    height: 100%;
    background: linear-gradient(90deg, #03dac6, #bb86fc);
    border-radius: 12px;
    transition: width 0.3s ease;
  }
  
  .progress-value {
    position: absolute;
    right: 10px;
    top: 50%;
    transform: translateY(-50%);
    font-size: 0.9rem;
    font-weight: 600;
    color: #fff;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
  }
  
  /* Rank badge styles */
  .rank-badge {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 36px;
    height: 36px;
    border-radius: 50%;
    color: white;
    font-weight: bold;
    font-size: 0.9rem;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
    margin-left: 8px;
    flex-shrink: 0;
    border: 2px solid rgba(255, 255, 255, 0.1);
  }
  
  .rank-gold {
    background: linear-gradient(135deg, #f9d423 0%, #f83600 100%);
  }
  
  .rank-silver {
    background: linear-gradient(135deg, #e6e6e6 0%, #8e8e8e 100%);
  }
  
  .rank-bronze {
    background: linear-gradient(135deg, #cd7f32 0%, #a05a2c 100%);
  }
  
  .rank-default {
    background: #424242;
  }
  
  /* Text container - typing area */
  .text-container {
    width: 100%;
    max-width: 800px;
    margin: 0 auto;
    padding: 2rem;
    border: 1px solid #333;
    border-radius: 12px;
    text-align: left;
    line-height: 1.8;
    min-height: 200px;
    position: relative;
    cursor: text;
    word-spacing: 5px;
    background-color: #1e1e1e;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2), inset 0 1px 2px rgba(255, 255, 255, 0.05);
    font-size: 1.2rem;
    transition: all 0.3s ease;
  }
  
  .text-container:focus {
    outline: none;
    border-color: #bb86fc;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2), 0 0 0 2px rgba(187, 134, 252, 0.3);
  }
  
  .text-container.rtl {
    direction: rtl;
    text-align: right;
  }
  
  .text-container.error-state {
    background-color: rgba(255, 82, 82, 0.1);
    border-color: #ff5252;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2), 0 0 0 2px rgba(255, 82, 82, 0.3);
  }
  
  /* Typing characters styles */
  :global(.word) {
    display: inline-block;
    margin: 0 2px;
    white-space: nowrap;
    border-radius: 3px;
    padding: 0 1px;
    transition: all 0.15s ease;
  }
  
  :global(.correct) {
    color: #03dac6; /* Teal accent */
  }
  
  :global(.incorrect) {
    color: #ff5252; /* Red accent */
    text-decoration: underline;
    text-decoration-color: rgba(255, 82, 82, 0.7);
    text-decoration-thickness: 2px;
  }
  
  :global(.pending) {
    color: #757575; /* Light gray */
  }
  
  :global(.current) {
    position: relative;
  }
  
  :global(.current-char) {
    background-color: rgba(187, 134, 252, 0.3); /* Purple with transparency */
    border-radius: 3px;
    animation: pulse 1.5s infinite;
  }
  
  :global(.current-char.error) {
    background-color: rgba(255, 82, 82, 0.3); /* Red with transparency */
  }
  
  :global(.error) {
    background-color: rgba(255, 82, 82, 0.3) !important;
    color: #ff5252 !important;
    border-bottom: 2px solid #ff5252;
  }
  
  /* Countdown timer */
  .countdown-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    width: 100%;
    max-width: 800px;
    min-height: 300px;
    background-color: #1e1e1e;
    border-radius: 12px;
    padding: 2rem;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
    margin-bottom: 2rem;
    animation: gentle-float 3s infinite ease-in-out;
    border: 1px solid #333;
  }
  
  @keyframes gentle-float {
    0%, 100% { transform: translateY(0); }
    50% { transform: translateY(-10px); }
  }
  
  .countdown-timer {
    font-size: 5rem;
    font-weight: bold;
    color: #bb86fc;
    text-shadow: 0 2px 8px rgba(187, 134, 252, 0.5);
    margin: 1rem 0;
    font-family: monospace;
    letter-spacing: -2px;
  }
  
  /* Waiting screen */
  .waiting {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    max-width: 800px;
    min-height: 200px;
    background-color: #1e1e1e;
    padding: 2rem;
    border-radius: 12px;
    font-size: 1.5rem;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
    border: 1px solid #333;
  }
  
  .waiting p {
    color: #bb86fc;
    position: relative;
    padding-right: 20px;
  }
  
  .waiting p::after {
    content: '';
    position: absolute;
    right: 0;
    top: 50%;
    width: 10px;
    height: 10px;
    background-color: #bb86fc;
    border-radius: 50%;
    transform: translateY(-50%);
    animation: blink 1s infinite;
  }
  
  /* Game End Popup Styles */
  .game-end-popup {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.85);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
    animation: fadeIn 0.3s ease;
    backdrop-filter: blur(5px);
  }
  
  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }
  
  .popup-content {
    background-color: #1e1e1e;
    padding: 2.5rem;
    border-radius: 16px;
    box-shadow: 0 16px 32px rgba(0, 0, 0, 0.4);
    width: 90%;
    max-width: 550px;
    text-align: center;
    animation: popIn 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
    max-height: 90vh;
    overflow-y: auto;
    border: 1px solid #333;
  }
  
  @keyframes popIn {
    from { transform: scale(0.9); opacity: 0; }
    to { transform: scale(1); opacity: 1; }
  }
  
  .popup-content h3 {
    color: #bb86fc;
    font-size: 2.5rem;
    margin-bottom: 1.5rem;
    font-weight: 600;
  }
  
  .results {
    margin-bottom: 2rem;
    text-align: left;
    background-color: #252525;
    padding: 1.5rem;
    border-radius: 12px;
    border: 1px solid #333;
  }
  
  .results p {
    font-weight: 600;
    margin-bottom: 0.5rem;
    color: #e0e0e0;
  }
  
  .results ul {
    list-style-type: none;
    padding: 0;
  }
  
  .results li {
    padding: 0.8rem 0;
    border-bottom: 1px solid #333;
    display: flex;
    justify-content: space-between;
  }
  
  .results li:last-child {
    border-bottom: none;
  }
  
  /* Final Leaderboard Styles */
  .final-leaderboard {
    margin-bottom: 2rem;
    background-color: #252525;
    padding: 1.5rem;
    border-radius: 12px;
    border: 1px solid #333;
  }
  
  .final-leaderboard table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 1rem;
  }
  
  .final-leaderboard th {
    text-align: left;
    padding: 0.75rem;
    border-bottom: 2px solid #333;
    font-weight: 600;
    color: #bb86fc;
  }
  
  .final-leaderboard td {
    padding: 1rem 0.75rem;
    border-bottom: 1px solid #333;
  }
  
  .final-leaderboard tr.current-user {
    background-color: rgba(187, 134, 252, 0.15);
    font-weight: 600;
  }
  
  .final-leaderboard tr:last-child td {
    border-bottom: none;
  }

  .final-leaderboard th:nth-child(3),
  .final-leaderboard th:nth-child(4) {
      text-align: center;
      width: 80px;
  }

  .wpm-cell, .duration-cell {
    text-align: center;
    font-weight: 600;
    color: #03dac6;
  }
  
  .popup-buttons {
    display: flex;
    justify-content: center;
    gap: 1.5rem;
    margin-top: 2rem;
  }
  
  .btn-main, .btn-play {
    padding: 0.75rem 1.75rem;
    border: none;
    border-radius: 8px;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.2);
  }
  
  .btn-main {
    background-color: #2d2d2d;
    color: #e0e0e0;
    border: 1px solid #424242;
  }
  
  .btn-main:hover {
    background-color: #3d3d3d;
    transform: translateY(-2px);
    box-shadow: 0 6px 10px rgba(0, 0, 0, 0.3);
  }
  
  .btn-play {
    background-color: #bb86fc;
    color: #121212;
  }
  
  .btn-play:hover {
    background-color: #c9a0ff;
    transform: translateY(-2px);
    box-shadow: 0 6px 10px rgba(0, 0, 0, 0.3);
  }
  
  @keyframes blink {
    0%, 100% { opacity: 1; }
    50% { opacity: 0; }
  }
  
  /* Responsive adjustments */
  @media (max-width: 768px) {
    .stats {
      flex-wrap: wrap;
    }
    
    .stat {
      min-width: 45%;
    }
    
    .countdown-timer {
      font-size: 3.5rem;
    }
    
    .text-container {
      font-size: 1rem;
      padding: 1.5rem;
    }
  }
  
  @media (max-width: 480px) {
    .stat {
      min-width: 100%;
    }
    
    .player-name {
      flex: 0 0 80px;
      font-size: 0.9rem;
    }
    
    .popup-content {
      padding: 1.5rem;
    }
  }
</style>
