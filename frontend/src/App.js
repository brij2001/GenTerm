import React, { useState, useEffect, useRef } from 'react';
import FileUploader from './components/FileUploader';
import sessionService from './services/sessionService';
import chatService from './services/chatService';
import fileService from './services/fileService';
import './App.css';

function App() {
  const [sessionId, setSessionId] = useState(null);
  const [isProcessing, setIsProcessing] = useState(false);
  const [uploadedFiles, setUploadedFiles] = useState([]);
  const [terminalHistory, setTerminalHistory] = useState([]);
  const [currentInput, setCurrentInput] = useState('');
  const [commandHistory, setCommandHistory] = useState([]);
  const [historyIndex, setHistoryIndex] = useState(-1);
  const inputRef = useRef(null);
  const terminalRef = useRef(null);
  const initializedRef = useRef(false);

  useEffect(() => {
    // Prevent double initialization due to StrictMode
    if (initializedRef.current) return;
    initializedRef.current = true;
    
    // Initialize terminal with welcome message
    addToTerminal('Welcome to GenTerm - Terminal-based RAG Q&A', 'system');
    addToTerminal('----------------------------------------', 'system');
    addToTerminal('Upload files with the panel on the right.', 'system');
    addToTerminal('Type your question and press Enter to ask.', 'system');
    addToTerminal('Type "clear" to clear the terminal.', 'system');
    addToTerminal('Type "files" to see uploaded files.', 'system');
    addToTerminal('Type "help" to display this list of commands.', 'system');
    
    // Create session
    createSession();
  }, []);
  
  // Focus input when component renders and when processing state changes
  useEffect(() => {
    if (inputRef.current) {
      setTimeout(() => {
        inputRef.current.focus();
      }, 100);
    }
  }, [isProcessing]);
  
  useEffect(() => {
    // Scroll to bottom whenever terminal history changes
    if (terminalRef.current) {
      terminalRef.current.scrollTop = terminalRef.current.scrollHeight;
    }
    
    // Focus input after terminal updates
    if (inputRef.current) {
      inputRef.current.focus();
    }
  }, [terminalHistory]);

  const createSession = async () => {
    try {
      const id = await sessionService.createSession();
      setSessionId(id);
      addToTerminal(`Session connected: ${id}`, 'system');
    } catch (error) {
      console.error('Failed to create session:', error);
      addToTerminal('Failed to create session. Please refresh the page.', 'error');
    }
  };

  const addToTerminal = (text, type = 'system') => {
    setTerminalHistory(prev => [...prev, { text, type }]);
  };

  const handleInputChange = (e) => {
    setCurrentInput(e.target.value);
  };

  const handleKeyDown = (e) => {
    if (isProcessing) return;

    if (e.key === 'Enter') {
      handleCommand();
    } else if (e.key === 'ArrowUp') {
      navigateHistory('up');
      e.preventDefault();
    } else if (e.key === 'ArrowDown') {
      navigateHistory('down');
      e.preventDefault();
    }
  };

  const navigateHistory = (direction) => {
    if (commandHistory.length === 0) return;

    if (direction === 'up' && historyIndex < commandHistory.length - 1) {
      const newIndex = historyIndex + 1;
      setHistoryIndex(newIndex);
      setCurrentInput(commandHistory[commandHistory.length - 1 - newIndex]);
    } else if (direction === 'down' && historyIndex > -1) {
      const newIndex = historyIndex - 1;
      setHistoryIndex(newIndex);
      if (newIndex === -1) {
        setCurrentInput('');
      } else {
        setCurrentInput(commandHistory[commandHistory.length - 1 - newIndex]);
      }
    }
  };

  const handleCommand = async () => {
    const command = currentInput.trim();
    
    // Display the command
    if (command) {
      addToTerminal(`${command}`, 'user');
    }
    
    // Add to history if not empty
    if (command) {
      setCommandHistory(prev => [...prev, command]);
      setHistoryIndex(-1);
    }

    setCurrentInput('');

    if (!command) {
      return;
    }

    switch (command) {
      case 'clear':
        setTerminalHistory([]);
        break;
      case 'help':
        showHelp();
        break;
      case 'files':
        listFiles();
        break;
      default:
        await processQuery(command);
        break;
    }
  };

  const showHelp = () => {
    addToTerminal('Available Commands:', 'system');
    addToTerminal('clear - Clear the terminal', 'system');
    addToTerminal('files - List uploaded files', 'system');
    addToTerminal('Type "help" to display this list of commands.', 'system');
    addToTerminal('', 'system');
    addToTerminal('Any other input will be treated as a question for the AI.', 'system');
  };

  const listFiles = () => {
    if (uploadedFiles.length === 0) {
      addToTerminal('No files uploaded.', 'system');
      return;
    }

    addToTerminal('Uploaded Files:', 'system');
    uploadedFiles.forEach((file, index) => {
      addToTerminal(`${index + 1}. ${file.name} (${fileService.formatFileSize(file.size)})`, 'system');
    });
  };

  const processQuery = async (query) => {
    if (!sessionId) {
      addToTerminal('No active session. Please refresh the page.', 'error');
      return;
    }

    setIsProcessing(true);
    addToTerminal(`Processing query... ${uploadedFiles.length} files in context.`, 'system');

    try {
      // Extract context from files
      const contexts = await Promise.all(
        uploadedFiles.map(file => fileService.extractContent(file))
      );

      const flattenedContext = contexts
        .filter(content => content) // Remove nulls
        .map((content, index) => `[File: ${uploadedFiles[index].name}]\n${content}`)
        .join('\n\n');

      // Check if there's an image selected for this query
      const imageFiles = uploadedFiles.filter(file => 
        file.type.includes('image/') || file.name.match(/\.(jpg|jpeg|png)$/i)
      );
      
      let response;
      
      if (imageFiles.length > 0 && query.toLowerCase().includes('image')) {
        // Get the last uploaded image file
        const imageFile = imageFiles[imageFiles.length - 1];
        addToTerminal(`Processing image: ${imageFile.name}...`, 'system');
        
        // Convert image to base64
        const base64Image = await fileService.imageToBase64(imageFile);
        
        // Get AI response with image
        addToTerminal('Analyzing image...', 'system');
        response = await chatService.sendImageQuery(sessionId, query, base64Image, flattenedContext);
      } else {
        // Regular text query
        addToTerminal('Thinking...', 'system');
        response = await chatService.sendQuery(sessionId, query, flattenedContext);
      }
      
      // Display response
      addToTerminal('AI Response:', 'system');
      addToTerminal('------------', 'system');
      addToTerminal(response, 'assistant');
      
    } catch (error) {
      console.error('Error processing query:', error);
      addToTerminal(`Error: ${error.message || 'Failed to process query'}`, 'error');
    } finally {
      setIsProcessing(false);
    }
  };

  const handleFilesUploaded = (files) => {
    setUploadedFiles(prevFiles => [...prevFiles, ...files]);
    
    addToTerminal(`Uploaded ${files.length} file(s):`, 'system');
    files.forEach(file => {
      addToTerminal(` - ${file.name} (${fileService.formatFileSize(file.size)})`, 'system');
    });
  };

  const getTerminalLineClass = (type) => {
    switch (type) {
      case 'user':
        return 'terminal-user';
      case 'assistant':
        return 'terminal-assistant';
      case 'error':
        return 'terminal-error';
      default:
        return 'terminal-system';
    }
  };

  const getPrompt = () => {
    return isProcessing ? 'processing... ' : 'genterm> ';
  };

  // Handle click on terminal to focus input
  const handleTerminalClick = () => {
    if (inputRef.current) {
      inputRef.current.focus();
    }
  };

  return (
    <div className="app">
      <header className="app-header">
        <h1>GenTerm</h1>
        <p>Terminal-based RAG Q&A Assistant</p>
      </header>
      <div className="app-container">
        <div className="terminal-container">
          <div 
            className="terminal" 
            ref={terminalRef}
            onClick={handleTerminalClick}
          >
            {terminalHistory.map((line, i) => (
              <div key={i} className={`terminal-line ${getTerminalLineClass(line.type)}`}>
                {line.type === 'user' && <span className="terminal-prompt">genterm&gt; </span>}
                <span className="terminal-text">{line.text}</span>
              </div>
            ))}
            <div className="terminal-input-line">
              <span className="terminal-prompt">{getPrompt()}</span>
              <input
                ref={inputRef}
                type="text"
                className="terminal-input"
                value={currentInput}
                onChange={handleInputChange}
                onKeyDown={handleKeyDown}
                disabled={isProcessing}
                autoFocus
              />
            </div>
          </div>
        </div>
        <div className="sidebar">
          <FileUploader onFilesUploaded={handleFilesUploaded} />
          <div className="file-list">
            <h3>Uploaded Files ({uploadedFiles.length})</h3>
            {uploadedFiles.length === 0 ? (
              <p>No files uploaded</p>
            ) : (
              <ul>
                {uploadedFiles.map((file, index) => (
                  <li key={index}>
                    {file.name} ({fileService.formatFileSize(file.size)})
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default App; 
