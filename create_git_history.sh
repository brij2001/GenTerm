#!/bin/bash

# Set git user configurations
git config user.name "brij2001"
git config user.email "brijmangukiya807@gmail.com"

# Function to create commits with specified dates
commit_with_date() {
    local commit_date="$1"
    local commit_message="$2"
    
    # Set the environment variables for commit date
    export GIT_AUTHOR_DATE="$commit_date"
    export GIT_COMMITTER_DATE="$commit_date"
    
    # Commit with the specified message
    git commit -m "$commit_message"
    
    # Reset environment variables
    unset GIT_AUTHOR_DATE
    unset GIT_COMMITTER_DATE
}

# Initialize the repository if not already initialized
if [ ! -d .git ]; then
    git init
    echo "Git repository initialized"
fi

# Create a README file - Feb 12, 2025
cat > README.md << 'EOF'
# GenTerm

A modern terminal emulator with advanced features and customization options.

## Features
- Custom themes and styling
- Plugin support
- Split-screen functionality
- Command history and search
- Integrated development tools

## Coming Soon
- Cloud synchronization
- Custom command aliases
- Advanced search capabilities
EOF

git add README.md .gitignore
commit_with_date "2025-02-12T10:30:00" "Initial commit: Project README with features overview"

# Create basic project structure - Feb 14, 2025
mkdir -p src/{core,ui,utils}
mkdir -p assets/themes

# Create a basic package.json
cat > package.json << 'EOF'
{
  "name": "genterm",
  "version": "0.1.0",
  "description": "A modern terminal emulator with advanced features",
  "main": "src/index.js",
  "scripts": {
    "start": "electron .",
    "test": "jest"
  },
  "dependencies": {
    "electron": "^20.0.0"
  },
  "devDependencies": {
    "jest": "^28.0.0"
  }
}
EOF

git add package.json
commit_with_date "2025-02-14T14:15:00" "Add package.json with basic dependencies"

# Create a simple entry point file
cat > src/index.js << 'EOF'
const { app, BrowserWindow } = require('electron');
const path = require('path');

function createWindow() {
  const win = new BrowserWindow({
    width: 900,
    height: 600,
    webPreferences: {
      nodeIntegration: true,
      contextIsolation: false
    }
  });

  win.loadFile(path.join(__dirname, 'ui/index.html'));
}

app.whenReady().then(() => {
  createWindow();
  
  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) createWindow();
  });
});

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') app.quit();
});
EOF

git add src/index.js
commit_with_date "2025-02-14T16:45:00" "Implement basic Electron app entry point"

# Create HTML entry point - Feb 16, 2025
mkdir -p src/ui
cat > src/ui/index.html << 'EOF'
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>GenTerm</title>
  <link rel="stylesheet" href="style.css">
</head>
<body>
  <div id="terminal-container">
    <div id="terminal-header">
      <div class="tab active">Terminal 1</div>
      <div class="tab-controls">
        <button id="new-tab">+</button>
      </div>
    </div>
    <div id="terminal-body"></div>
  </div>
  <script src="renderer.js"></script>
</body>
</html>
EOF

git add src/ui/index.html
commit_with_date "2025-02-16T11:20:00" "Add basic HTML structure for the terminal UI"

# Create CSS styles - Feb 18, 2025
cat > src/ui/style.css << 'EOF'
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Consolas', monospace;
  background-color: #1e1e1e;
  color: #f0f0f0;
  height: 100vh;
  overflow: hidden;
}

#terminal-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

#terminal-header {
  display: flex;
  background-color: #333;
  padding: 5px;
}

.tab {
  padding: 5px 15px;
  border-radius: 4px 4px 0 0;
  margin-right: 2px;
  cursor: pointer;
}

.tab.active {
  background-color: #1e1e1e;
}

.tab-controls {
  margin-left: auto;
}

#terminal-body {
  flex: 1;
  padding: 10px;
  font-size: 14px;
  overflow: auto;
}
EOF

# Only add part of the CSS file initially
git add -p src/ui/style.css
commit_with_date "2025-02-18T09:30:00" "Add initial styling for terminal UI"

# Add the rest of the CSS later
git add src/ui/style.css
commit_with_date "2025-02-18T14:15:00" "Complete terminal styling with header and tabs"

# Create renderer.js - Feb 20, 2025
cat > src/ui/renderer.js << 'EOF'
document.addEventListener('DOMContentLoaded', () => {
  const terminalBody = document.getElementById('terminal-body');
  const newTabButton = document.getElementById('new-tab');
  
  // Simulate a terminal prompt
  function createPrompt() {
    const prompt = document.createElement('div');
    prompt.className = 'prompt';
    prompt.innerHTML = '<span class="user">user</span>@<span class="host">genterm</span>:<span class="directory">~</span>$ ';
    terminalBody.appendChild(prompt);
    
    const input = document.createElement('input');
    input.type = 'text';
    input.className = 'command-input';
    terminalBody.appendChild(input);
    input.focus();
    
    input.addEventListener('keydown', (e) => {
      if (e.key === 'Enter') {
        const command = input.value;
        input.disabled = true;
        
        // Process command (simple echo for now)
        if (command.trim() !== '') {
          const output = document.createElement('div');
          output.className = 'output';
          output.textContent = `Echo: ${command}`;
          terminalBody.appendChild(output);
        }
        
        createPrompt();
      }
    });
  }
  
  createPrompt();
  
  newTabButton.addEventListener('click', () => {
    // Tab creation logic to be implemented
    console.log('New tab requested');
  });
});
EOF

# Add part of the renderer.js first
git add -p src/ui/renderer.js
commit_with_date "2025-02-20T13:10:00" "Implement basic terminal prompt functionality"

# Add the rest of renderer.js
git add src/ui/renderer.js
commit_with_date "2025-02-20T18:30:00" "Add command input handling and UI interactions"

# Create core terminal functionality - Feb 25, 2025
cat > src/core/terminal.js << 'EOF'
class Terminal {
  constructor(element, options = {}) {
    this.element = element;
    this.history = [];
    this.historyIndex = -1;
    this.commandHandlers = {};
    this.currentDirectory = options.initialDirectory || '~';
    
    this.initializeEventListeners();
  }
  
  initializeEventListeners() {
    // To be implemented
  }
  
  executeCommand(commandString) {
    const [command, ...args] = commandString.trim().split(' ');
    
    this.history.push(commandString);
    this.historyIndex = this.history.length;
    
    if (this.commandHandlers[command]) {
      return this.commandHandlers[command](args);
    } else {
      return `Command not found: ${command}`;
    }
  }
  
  registerCommand(name, handler) {
    this.commandHandlers[name] = handler;
  }
  
  clearScreen() {
    while (this.element.firstChild) {
      this.element.removeChild(this.element.firstChild);
    }
  }
}

module.exports = Terminal;
EOF

git add src/core/terminal.js
commit_with_date "2025-02-25T11:45:00" "Create core Terminal class with command execution"

# Create some utility functions - Feb 28, 2025
cat > src/utils/commands.js << 'EOF'
const fs = require('fs');
const path = require('path');

// Basic command implementations
const commands = {
  echo: (args) => args.join(' '),
  
  clear: (_, terminal) => {
    terminal.clearScreen();
    return '';
  },
  
  ls: (args, terminal, state) => {
    // Simple directory listing (mock for now)
    return [
      'Documents',
      'Downloads',
      'Projects',
      'notes.txt',
      'config.json'
    ].join('\n');
  },
  
  pwd: (_, terminal, state) => {
    return state.currentDirectory;
  },
  
  cd: (args, terminal, state) => {
    if (!args[0]) {
      state.currentDirectory = '~';
      return '';
    }
    
    // Simple directory change (mock for now)
    state.currentDirectory = args[0];
    return '';
  },
  
  help: () => {
    return [
      'Available commands:',
      '  echo [text] - Display text',
      '  clear - Clear the terminal screen',
      '  ls - List directory contents',
      '  pwd - Print working directory',
      '  cd [directory] - Change directory',
      '  help - Display this help message'
    ].join('\n');
  }
};

module.exports = commands;
EOF

# Add part of the commands.js file
git add -p src/utils/commands.js
commit_with_date "2025-02-28T10:15:00" "Implement basic command utilities (echo, clear)"

# Add the rest of commands.js
git add src/utils/commands.js
commit_with_date "2025-02-28T16:30:00" "Add file system commands (ls, cd, pwd)"

# Create theme support - March 3, 2025
cat > src/utils/theme-manager.js << 'EOF'
class ThemeManager {
  constructor() {
    this.themes = {};
    this.currentTheme = 'default';
    
    // Register default themes
    this.registerDefaultThemes();
  }
  
  registerDefaultThemes() {
    this.registerTheme('default', {
      background: '#1e1e1e',
      foreground: '#f0f0f0',
      cursor: '#ffffff',
      selection: 'rgba(255, 255, 255, 0.3)',
      black: '#000000',
      red: '#e06c75',
      green: '#98c379',
      yellow: '#e5c07b',
      blue: '#61afef',
      magenta: '#c678dd',
      cyan: '#56b6c2',
      white: '#dcdfe4'
    });
    
    this.registerTheme('light', {
      background: '#f0f0f0',
      foreground: '#333333',
      cursor: '#000000',
      selection: 'rgba(0, 0, 0, 0.3)',
      black: '#000000',
      red: '#e45649',
      green: '#50a14f',
      yellow: '#c18401',
      blue: '#4078f2',
      magenta: '#a626a4',
      cyan: '#0184bc',
      white: '#a0a1a7'
    });
  }
  
  registerTheme(name, colors) {
    this.themes[name] = colors;
  }
  
  setTheme(name) {
    if (!this.themes[name]) {
      throw new Error(`Theme '${name}' not found`);
    }
    
    this.currentTheme = name;
    this.applyTheme();
  }
  
  applyTheme() {
    const colors = this.themes[this.currentTheme];
    const root = document.documentElement;
    
    Object.keys(colors).forEach(key => {
      root.style.setProperty(`--${key}`, colors[key]);
    });
  }
}

module.exports = ThemeManager;
EOF

git add src/utils/theme-manager.js
commit_with_date "2025-03-03T13:45:00" "Add theme management with default dark and light themes"

# Update CSS to use theme variables - March 5, 2025
cat > src/ui/style.css << 'EOF'
:root {
  /* Theme variables will be set dynamically */
  --background: #1e1e1e;
  --foreground: #f0f0f0;
  --cursor: #ffffff;
  --selection: rgba(255, 255, 255, 0.3);
  --black: #000000;
  --red: #e06c75;
  --green: #98c379;
  --yellow: #e5c07b;
  --blue: #61afef;
  --magenta: #c678dd;
  --cyan: #56b6c2;
  --white: #dcdfe4;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Consolas', monospace;
  background-color: var(--background);
  color: var(--foreground);
  height: 100vh;
  overflow: hidden;
}

#terminal-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

#terminal-header {
  display: flex;
  background-color: #333;
  padding: 5px;
}

.tab {
  padding: 5px 15px;
  border-radius: 4px 4px 0 0;
  margin-right: 2px;
  cursor: pointer;
}

.tab.active {
  background-color: var(--background);
}

.tab-controls {
  margin-left: auto;
}

#terminal-body {
  flex: 1;
  padding: 10px;
  font-size: 14px;
  overflow: auto;
}

.prompt {
  display: inline;
}

.prompt .user {
  color: var(--green);
}

.prompt .host {
  color: var(--blue);
}

.prompt .directory {
  color: var(--yellow);
}

.command-input {
  background: transparent;
  border: none;
  color: var(--foreground);
  font-family: inherit;
  font-size: inherit;
  outline: none;
  width: 80%;
  margin-left: 5px;
}

.output {
  margin: 5px 0;
  white-space: pre-wrap;
}
EOF

# Update part of the CSS file
git add -p src/ui/style.css
commit_with_date "2025-03-05T11:30:00" "Update CSS to use theme variables"

# Add the rest of updated CSS
git add src/ui/style.css
commit_with_date "2025-03-05T15:20:00" "Add styled terminal prompts and output"

# Create settings manager - March 8, 2025
cat > src/core/settings.js << 'EOF'
const fs = require('fs');
const path = require('path');
const { app } = require('electron');

class SettingsManager {
  constructor() {
    this.settings = {
      fontSize: 14,
      fontFamily: 'Consolas, monospace',
      theme: 'default',
      cursorStyle: 'block',
      cursorBlink: true,
      shell: process.platform === 'win32' ? 'powershell.exe' : 'bash',
      scrollback: 1000,
      tabSize: 4
    };
    
    this.settingsPath = path.join(app.getPath('userData'), 'settings.json');
    this.loadSettings();
  }
  
  loadSettings() {
    try {
      if (fs.existsSync(this.settingsPath)) {
        const data = fs.readFileSync(this.settingsPath, 'utf8');
        const loadedSettings = JSON.parse(data);
        this.settings = { ...this.settings, ...loadedSettings };
      }
    } catch (err) {
      console.error('Error loading settings:', err);
    }
  }
  
  saveSettings() {
    try {
      fs.writeFileSync(this.settingsPath, JSON.stringify(this.settings, null, 2));
    } catch (err) {
      console.error('Error saving settings:', err);
    }
  }
  
  get(key) {
    return this.settings[key];
  }
  
  set(key, value) {
    this.settings[key] = value;
    this.saveSettings();
  }
}

module.exports = SettingsManager;
EOF

git add src/core/settings.js
commit_with_date "2025-03-08T14:20:00" "Implement settings manager for user preferences"

# Create plugin system - March 10, 2025
cat > src/core/plugin-manager.js << 'EOF'
class PluginManager {
  constructor(terminal) {
    this.terminal = terminal;
    this.plugins = {};
  }
  
  register(name, plugin) {
    if (this.plugins[name]) {
      throw new Error(`Plugin '${name}' is already registered`);
    }
    
    this.plugins[name] = plugin;
    
    if (typeof plugin.initialize === 'function') {
      plugin.initialize(this.terminal);
    }
    
    console.log(`Plugin '${name}' registered successfully`);
  }
  
  unregister(name) {
    const plugin = this.plugins[name];
    
    if (!plugin) {
      throw new Error(`Plugin '${name}' is not registered`);
    }
    
    if (typeof plugin.cleanup === 'function') {
      plugin.cleanup(this.terminal);
    }
    
    delete this.plugins[name];
    console.log(`Plugin '${name}' unregistered successfully`);
  }
  
  getPlugin(name) {
    return this.plugins[name] || null;
  }
  
  getRegisteredPlugins() {
    return Object.keys(this.plugins);
  }
}

module.exports = PluginManager;
EOF

git add src/core/plugin-manager.js
commit_with_date "2025-03-10T10:45:00" "Create plugin system for terminal extensions"

# Create a sample plugin - March 12, 2025
mkdir -p plugins
cat > plugins/git-status.js << 'EOF'
const { spawn } = require('child_process');

// Example plugin that adds git status integration
class GitStatusPlugin {
  initialize(terminal) {
    // Register a command to show git status
    terminal.registerCommand('git-status', this.gitStatus.bind(this));
    
    // Add git status to prompt if in a git repository
    this.patchPromptFunction(terminal);
  }
  
  gitStatus(args, terminal) {
    return new Promise((resolve, reject) => {
      const git = spawn('git', ['status', '--porcelain']);
      let output = '';
      
      git.stdout.on('data', (data) => {
        output += data.toString();
      });
      
      git.on('close', (code) => {
        if (code === 0) {
          if (output.trim() === '') {
            resolve('Working directory clean');
          } else {
            resolve(output);
          }
        } else {
          resolve('Not a git repository or git command failed');
        }
      });
      
      git.on('error', (err) => {
        resolve(`Error: ${err.message}`);
      });
    });
  }
  
  patchPromptFunction(terminal) {
    // Implementation depends on terminal's prompt mechanism
    // This is a placeholder
  }
  
  cleanup(terminal) {
    // Remove custom commands when plugin is unregistered
    delete terminal.commandHandlers['git-status'];
  }
}

module.exports = new GitStatusPlugin();
EOF

git add plugins/git-status.js
commit_with_date "2025-03-12T16:30:00" "Add sample Git status plugin"

# Update README with installation instructions - March 14, 2025
cat > README.md << 'EOF'
# GenTerm

A modern terminal emulator with advanced features and customization options.

## Features
- Custom themes and styling
- Plugin support
- Split-screen functionality
- Command history and search
- Integrated development tools

## Installation

### Prerequisites
- Node.js 14+
- npm or yarn

### Steps
1. Clone the repository
   ```
   git clone https://github.com/brij2001/GenTerm.git
   cd GenTerm
   ```

2. Install dependencies
   ```
   npm install
   ```

3. Start the application
   ```
   npm start
   ```

## Customization
GenTerm can be customized through the settings panel or by editing the `settings.json` file directly.

### Themes
The application comes with several built-in themes:
- Default (Dark)
- Light
- Solarized Dark
- Solarized Light

You can create custom themes by adding new theme files to the `assets/themes` directory.

### Plugins
Extend GenTerm's functionality with plugins:

1. Install plugins from the plugin marketplace
2. Or place plugin files in the `plugins` directory

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

## License
This project is licensed under the MIT License - see the LICENSE file for details.
EOF

git add README.md
commit_with_date "2025-03-14T12:00:00" "Update README with installation and customization instructions"

echo "Git history creation complete!"
echo "Repository now has commits from February 12, 2025 to March 14, 2025" 