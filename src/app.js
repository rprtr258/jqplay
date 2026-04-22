import * as jq from 'jq-wasm';

// Editor instances
let filterEditor, jsonEditor, resultEditor;

// Debounce timer
let runTimeout = null;

// Initialize editors
function initEditors() {
  // Filter editor
  filterEditor = ace.edit('filter-editor');
  filterEditor.setTheme('ace/theme/tomorrow_night');
  filterEditor.session.setMode('ace/mode/jsoniq');
  filterEditor.setHighlightActiveLine(false);
  filterEditor.setFontSize(14);
  filterEditor.setShowPrintMargin(false);
  filterEditor.session.setUseWorker(false);
  filterEditor.focus();

  // JSON editor
  jsonEditor = ace.edit('json-editor');
  jsonEditor.setTheme('ace/theme/tomorrow_night');
  jsonEditor.session.setMode('ace/mode/jsoniq');
  jsonEditor.setHighlightActiveLine(false);
  jsonEditor.setFontSize(14);
  jsonEditor.setShowPrintMargin(false);
  jsonEditor.session.setUseWorker(false);

  // Result editor (readonly)
  resultEditor = ace.edit('result-editor');
  resultEditor.setTheme('ace/theme/tomorrow_night');
  resultEditor.session.setMode('ace/mode/jsoniq');
  resultEditor.setHighlightActiveLine(false);
  resultEditor.setFontSize(14);
  resultEditor.setShowPrintMargin(false);
  resultEditor.session.setUseWorker(false);
  resultEditor.setReadOnly(true);

  // Set initial content
  filterEditor.setValue('. | with_entries({key: .key, value: .value.name})');
  jsonEditor.setValue(JSON.stringify({
    person1: {
      name: "Alice",
      welcome: "Hello Alice!"
    },
    person2: {
      name: "Bob",
      welcome: "Hello Bob!"
    }
  }, null, 2));
  resultEditor.setValue('');

  // Clear selection
  filterEditor.clearSelection();
  jsonEditor.clearSelection();
  resultEditor.clearSelection();

  // Add change listeners with debounce
  filterEditor.session.on('change', debounceRun);
  jsonEditor.session.on('change', debounceRun);

  // Add option change listeners
  document.querySelectorAll('.option-label input').forEach(input => {
    input.addEventListener('change', debounceRun);
  });

  // Run initial query
  run();
}

// Debounced run function
function debounceRun() {
  if (runTimeout) {
    clearTimeout(runTimeout);
  }
  runTimeout = setTimeout(run, 300);
}

// Collect flags from checkboxes
function getFlags() {
  const flags = [];

  document.querySelectorAll('.option-label input[type="checkbox"]:checked').forEach(input => {
    const flag = input.dataset.flag;
    if (flag) {
      flags.push(flag);
    }
  });

  // Handle indent flag (only if tab is not checked)
  const tabCheckbox = document.getElementById('opt-tab');
  const indentInput = document.getElementById('opt-indent');
  if (!tabCheckbox.checked && indentInput.value !== '2') {
    flags.push('--indent', indentInput.value);
  }

  return flags;
}

// Run jq query
async function run() {
  const query = filterEditor.getValue();
  const jsonInput = jsonEditor.getValue();
  const flags = getFlags();

  if (!query) {
    resultEditor.setValue('Error: missing filter');
    resultEditor.clearSelection();
    return;
  }

  if (!jsonInput) {
    resultEditor.setValue('Error: missing JSON');
    resultEditor.clearSelection();
    return;
  }

  resultEditor.setValue('Loading...');
  resultEditor.clearSelection();

  try {
    const result = await jq.raw(jsonInput, query, flags);

    if (result.stderr) {
      resultEditor.setValue(result.stderr);
    } else {
      resultEditor.setValue(result.stdout);
    }
  } catch (err) {
    resultEditor.setValue(`Error: ${err.message}`);
  }

  resultEditor.clearSelection();
}

// Initialize on DOM ready
document.addEventListener('DOMContentLoaded', initEditors);

// Resize handle functionality
document.addEventListener('DOMContentLoaded', () => {
  const handle = document.getElementById('filter-resize-handle');
  const filterEditorEl = document.getElementById('filter-editor');
  let isResizing = false;
  let startY = 0;
  let startHeight = 0;

  handle.addEventListener('mousedown', (e) => {
    isResizing = true;
    startY = e.clientY;
    startHeight = filterEditorEl.offsetHeight;
    handle.classList.add('active');
    document.body.style.cursor = 'row-resize';
    document.body.style.userSelect = 'none';
    e.preventDefault();
  });

  document.addEventListener('mousemove', (e) => {
    if (!isResizing) return;
    const deltaY = e.clientY - startY;
    const newHeight = Math.max(40, startHeight + deltaY);
    filterEditorEl.style.height = newHeight + 'px';
    // Trigger Ace editor resize
    filterEditor.resize();
  });

  document.addEventListener('mouseup', () => {
    if (isResizing) {
      isResizing = false;
      handle.classList.remove('active');
      document.body.style.cursor = '';
      document.body.style.userSelect = '';
    }
  });
});

// Cheatsheet modal
window.applyCheatsheetExample = (query, json) => {
  filterEditor.setValue(query);
  jsonEditor.setValue(json);
  filterEditor.clearSelection();
  jsonEditor.clearSelection();
  $('#cheatsheet-modal').modal('hide');
};

document.addEventListener('DOMContentLoaded', () => {
  const cheatsheetBtn = document.getElementById('cheatsheet-btn');
  if (cheatsheetBtn) {
    cheatsheetBtn.addEventListener('click', () => {
      $('#cheatsheet-modal').modal('show');
    });
  }

  // Handle cheatsheet row clicks
  document.querySelectorAll('.cheat-table tbody tr').forEach(row => {
    row.addEventListener('click', () => {
      const query = row.dataset.query;
      const json = row.dataset.json;
      if (query && json) {
        window.applyCheatsheetExample(query, json);
      }
    });
  });
});
