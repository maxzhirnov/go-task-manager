<script>
    import { browser } from '$app/environment';
    import { api } from '../api.js';
    import { tasks } from '../stores.js';
    import { mdiPlus, mdiChevronDown } from '@mdi/js';

    let title = '';
    let description = '';
    let status = 'pending';
    let isQuickAdd = true;


    async function handleQuickSubmit() {
        if (!title.trim()) return;
        
        if (title.length > 100) {
            showError("Task title cannot exceed 100 characters.");
            return;
        }

        await api.addTask(title, '', 'pending');
        const updatedTasks = await api.fetchTasks();
        tasks.set(updatedTasks);
        
        title = '';
    }

    async function handleFullSubmit() {
        if (title.length > 100) {
            showError("Task title cannot exceed 100 characters.");
            return;
        }

        await api.addTask(title, description, status);
        const updatedTasks = await api.fetchTasks();
        tasks.set(updatedTasks);
        
        title = '';
        description = '';
        status = 'pending';
        isQuickAdd = true;
    }
</script>

<!-- <div class="task-creation-container">
    {#if isQuickAdd}
        <div class="quick-add">
            <div class="input-wrapper">
                <input 
                    bind:value={title} 
                    type="text" 
                    placeholder="Quick add task" 
                    on:keydown={(e) => e.key === 'Enter' && handleQuickSubmit()}
                    maxlength="100"
                />
                <span class="quick-character-count" class:warning={title.length > 90}>
                    {title.length}/100
                </span>
            </div>
            <div class="quick-add-buttons">
                <div class="quick-submit-btn-container">
                    <button 
                    class="quick-submit-btn" 
                    title="Add task"
                    aria-label="Add task"
                    on:click={handleQuickSubmit}
                    disabled={!title.trim()}
                >
                    <svg viewBox="0 0 24 24" width="24" height="24">
                        <path fill="currentColor" d={mdiPlus} />
                    </svg>
                </button>
                </div>
                
                <div>
                    <button 
                        class="expand-btn" 
                        title="Show full form"
                        aria-label="Show full form"
                        on:click={() => isQuickAdd = false}
                    >
                        <svg viewBox="0 0 24 24" width="24" height="24">
                            <path fill="currentColor" d={mdiChevronDown} />
                        </svg>
                    </button>
                </div>
            </div>
        </div>
    {:else}
        <form class="task-form" on:submit|preventDefault={handleFullSubmit}>
            <div class="input-wrapper">
                <input 
                    bind:value={title} 
                    type="text" 
                    placeholder="Task Title" 
                    required
                    maxlength="100"
                    name="title"
                    on:keydown={(e) => e.key === 'Enter' && handleFullSubmit()}
                />
                <span class="quick-character-count" class:warning={title.length > 90}>
                    {title.length}/100
                </span>
            </div>

            <textarea 
                bind:value={description} 
                placeholder="Task Description" 
                name="description"
                on:keydown={(e) => e.key === 'Enter' && handleFullSubmit()}
            ></textarea>
            <select bind:value={status} required>
                <option value="pending">Pending</option>
                <option value="in_progress">In Progress</option>
                <option value="completed">Completed</option>
            </select>
            <div class="task-form-buttons">
                <button class="btn-add" type="submit">Add</button>
                <button 
                    class="btn-cancel" 
                    type="button" 
                    on:click={() => {
                        isQuickAdd = true;
                    }}
                >
                    Close
                </button>
            </div>
        </form>
    {/if}
</div> -->

<!-- <style>
    .input-wrapper {
        position: relative;
        flex: 1;
    }

    .input-wrapper input {
        width: 100%;
        padding-right: 60px; /* Make room for the counter */
    }

    .quick-character-count {
        position: absolute;
        right: 12px;
        top: 50%;
        transform: translateY(-50%);
        font-size: 12px;
        color: #666;
        transition: all 0.2s;
    }

    .quick-character-count.warning {
        color: #ff4444;
        font-weight: bold;
    }

    /* Container styles */
    .task-creation-container {
        margin: 3rem 0 1rem 0;
    }

    /* Common form element styles */
    .task-form input, .task-form textarea, .input-wrapper input {
        font-family: Arial, Helvetica, sans-serif;
        padding: 12px;
        border: 1px solid #ddd;
        border-radius: 4px;
        font-size: 16px;
        box-sizing: border-box;
        transition: border-color 0.2s, box-shadow 0.2s;
    }

    .task-form select {
        width: 100%;
        font-size: 16px;
        appearance: none;
        padding: 8px 32px 8px 12px;
        border: 1px solid #ddd;
        border-radius: 4px;
        background: #fff url("data:image/svg+xml;charset=US-ASCII,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%22292.4%22%20height%3D%22292.4%22%3E%3Cpath%20fill%3D%22%23666%22%20d%3D%22M287%2069.4a17.6%2017.6%200%200%200-13-5.4H18.4c-5%200-9.3%201.8-12.9%205.4A17.6%2017.6%200%200%200%200%2082.2c0%205%201.8%209.3%205.4%2012.9l128%20127.9c3.6%203.6%207.8%205.4%2012.8%205.4s9.2-1.8%2012.8-5.4L287%2095c3.5-3.5%205.4-7.8%205.4-12.8%200-5-1.9-9.2-5.4-12.8z%22%2F%3E%3C%2Fsvg%3E") no-repeat right 12px top 50%;
        background-size: 12px auto;
        font-size: 14px;
        cursor: pointer;
        transition: border-color 0.2s, box-shadow 0.2s;
        width: 100%;
    }

    .task-form select:hover {
        border-color: #aaa;
    }

    .task-form select:focus {
        outline: none;
        border-color: #2196F3;
        box-shadow: 0 0 0 2px rgba(33, 150, 243, 0.2);
    }

    .task-form select option {
        padding: 8px;
        background-color: white;
    }

    /* Quick add section */
    .quick-add {
        display: flex;
        gap: 8px;
        align-items: center;
    }

    .quick-add-buttons {
        gap: 8px;
        display: flex;
    }

    .quick-add input {
        flex: 1;
    }

    /* Action buttons */
    .quick-submit-btn,
    .expand-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        min-width: 42px;
        height: 42px;
        padding: 8px;
        border-radius: 4px;
        cursor: pointer;
        transition: all 0.2s;
    }

    .quick-submit-btn {
        background-color: #4CAF50;
        border: 1px solid #4CAF50;
        color: white;
    }

    .quick-submit-btn:hover {
        background-color: #45a049;
    }

    .quick-submit-btn:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    .expand-btn {
        background: none;
        border: 1px solid #ddd;
    }

    .expand-btn:hover {
        background: #f5f5f5;
        border-color: #aaa;
    }

    /* Full form styles */
    .task-form {
        display: flex;
        flex-direction: column;
        gap: 10px;
        margin-top: 20px;
    }

    .task-form-buttons {
        display: flex;
        justify-content: flex-end;
        gap: .5rem;
    }

    .task-form button {
        min-width: 100px;
        padding: 12px;
        border: none;
        border-radius: 4px;
        color: white;
        cursor: pointer;
    }

    .btn-add {
        background-color: #32cd32;
    }

    .btn-add:hover {
        background-color: #45a049;
    }

    .btn-cancel {
        background-color: #7f7f7f;
    }

    .btn-cancel:hover {
        background-color: #5d5c5c;
    }

    .character-count {
        font-size: 12px;
        color: #666;
        margin-left: 5px;
    }

    /* Mobile styles */
    @media screen and (max-width: 600px) {
        .input-wrapper {
            width: 100%;
        }

        .task-creation-container {
            margin-top: 0;
            padding: 0;
        }

        .quick-add {
            display: flex;
            flex-direction: column;
            width: 100%;
            margin: 0;
        }

        .quick-add input {
            width: 100%;
        }

        .quick-add-buttons {
            display: flex;
            justify-content: flex-end;
            width: 100%;
            gap: 8px;
        }

        .quick-submit-btn-container {
            flex: 1;
        }

        .quick-submit-btn-container button {
            width: 100%;
        }

        .quick-submit-btn,
        .expand-btn {
            height: 44.5px;
        }
    }
</style> -->

<div class="task-creation-container">
    {#if isQuickAdd}
        <div class="quick-add">
            <div class="terminal-header">
                <span class="terminal-dots">
                    <span class="dot"></span>
                    <span class="dot"></span>
                    <span class="dot"></span>
                </span>
                <span class="terminal-title">QUICK_TASK.exe</span>
            </div>
            <div class="input-wrapper">
                <span class="prompt">>_</span>
                <input 
                    bind:value={title} 
                    type="text" 
                    placeholder="Initialize new task..." 
                    on:keydown={(e) => e.key === 'Enter' && handleQuickSubmit()}
                    maxlength="100"
                />
                <span class="quick-character-count" class:warning={title.length > 90}>
                    [{title.length}/100]
                </span>
            </div>
            <div class="quick-add-buttons">
                <div class="quick-submit-btn-container">
                    <button 
                        class="quick-submit-btn" 
                        title="Execute task"
                        aria-label="Add task"
                        on:click={handleQuickSubmit}
                        disabled={!title.trim()}
                    >
                        <svg viewBox="0 0 24 24" width="24" height="24">
                            <path fill="currentColor" d={mdiPlus} />
                        </svg>
                        <span class="btn-label">EXECUTE</span>
                    </button>
                </div>
                
                <div>
                    <button 
                        class="expand-btn" 
                        title="Expand parameters"
                        aria-label="Show full form"
                        on:click={() => isQuickAdd = false}
                    >
                        <svg viewBox="0 0 24 24" width="24" height="24">
                            <path fill="currentColor" d={mdiChevronDown} />
                        </svg>
                        <span class="btn-label">CONFIG</span>
                    </button>
                </div>
            </div>
        </div>
    {:else}
        <form class="task-form" on:submit|preventDefault={handleFullSubmit}>
            <div class="terminal-header">
                <span class="terminal-dots">
                    <span class="dot"></span>
                    <span class="dot"></span>
                    <span class="dot"></span>
                </span>
                <span class="terminal-title">TASK_CONFIGURATION.exe</span>
            </div>
            
            <div class="form-grid">
                <div class="input-group">
                    <div class="input-label">[TITLE]</div>
                    <div class="input-wrapper">
                        <span class="prompt">>_</span>
                        <input 
                            bind:value={title} 
                            type="text" 
                            placeholder="Task identifier..." 
                            required
                            maxlength="100"
                            name="title"
                        />
                        <span class="quick-character-count" class:warning={title.length > 90}>
                            [{title.length}/100]
                        </span>
                    </div>
                </div>

                <div class="input-group">
                    <div class="input-label">[DESCRIPTION]</div>
                    <div class="input-wrapper">
                        <span class="prompt">>_</span>
                        <textarea 
                            bind:value={description} 
                            placeholder="Task parameters..." 
                            name="description"
                        ></textarea>
                    </div>
                </div>

                <div class="input-group">
                    <div class="input-label">[STATUS]</div>
                    <div class="select-wrapper">
                        <select bind:value={status} required>
                            <option value="pending">PENDING::0x01</option>
                            <option value="in_progress">IN_PROGRESS::0x02</option>
                            <option value="completed">COMPLETED::0x03</option>
                        </select>
                    </div>
                </div>
            </div>

            <div class="task-form-buttons">
                <button class="btn-cancel" type="button" on:click={() => isQuickAdd = true}>
                    <span class="btn-icon">×</span>
                    <span class="btn-label">ABORT</span>
                </button>
                <button class="btn-add" type="submit">
                    <span class="btn-icon">⚡</span>
                    <span class="btn-label">EXECUTE</span>
                </button>
            </div>

            <div class="form-status">
                <span class="status-indicator"></span>
                <span class="status-text">System Ready</span>
            </div>
        </form>
    {/if}
</div>

<style>
    .task-creation-container {
        background: #1c1c1c;
        border: 1px solid #0984e3;
        border-radius: 4px;
        padding: 1rem;
        margin-bottom: 1.5rem;
        position: relative;
        font-family: "JetBrains Mono", monospace;
        box-shadow: 0 2px 8px rgba(9, 132, 227, 0.1);
    }

    /* Circuit board pattern */
    .task-creation-container::before {
        content: "";
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background-image: 
            radial-gradient(
                circle at 50% 50%,
                rgba(0, 184, 148, 0.05) 1px,
                transparent 1px
            );
        background-size: 10px 10px;
        pointer-events: none;
    }

    .terminal-header {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        margin-bottom: 1rem;
        padding-bottom: 0.5rem;
        border-bottom: 1px solid rgba(9, 132, 227, 0.2);
    }

    .terminal-dots {
        display: flex;
        gap: 4px;
    }

    .dot {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background: #636e72;
    }

    .terminal-title {
        color: #00b894;
        font-size: 0.7rem;
        letter-spacing: 0.1em;
    }

    .quick-add {
        display: flex;
        gap: 1rem;
        align-items: center;
        position: relative;
        z-index: 1;
    }

    .input-wrapper {
        flex-grow: 1;
        position: relative;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        background: #2d3436;
        border: 1px solid #0984e3;
        border-radius: 3px;
        padding: 0 0.5rem;
    }

    .prompt {
        color: #00b894;
        font-size: 0.9rem;
        user-select: none;
    }

    input, textarea, select {
        width: 100%;
        background: transparent;
        border: none;
        color: #fff;
        padding: 0.8rem 3rem 0.8rem 0.5rem; /* Space for character count */
        font-family: "JetBrains Mono", monospace;
        font-size: 0.9rem;
    }

    input:focus, textarea:focus, select:focus {
        outline: none;
    }

    .input-wrapper:focus-within {
        border-color: #00b894;
        box-shadow: 0 0 10px rgba(0, 184, 148, 0.2);
    }

    .quick-character-count {
        position: absolute;
        right: 0.8rem;
        color: #636e72;
        font-size: 0.7rem;
        pointer-events: none;
    }

    .warning {
        color: #ff7675;
    }

    .quick-add-buttons {
        display: flex;
        gap: 0.5rem;
        flex-shrink: 0;
    }

    button {
        background: transparent;
        border: 1px solid #00b894;
        color: #00b894;
        padding: 0.5rem 1rem;
        border-radius: 3px;
        cursor: pointer;
        transition: all 0.3s ease;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-family: "JetBrains Mono", monospace;
        font-size: 0.8rem;
    }

    button:hover:not(:disabled) {
        background: rgba(0, 184, 148, 0.1);
        box-shadow: 0 0 8px rgba(0, 184, 148, 0.3);
    }

    button:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    /* Full form styles */
    .task-form {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .form-grid {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .input-group {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .input-label {
        color: #00b894;
        font-size: 0.7rem;
        letter-spacing: 0.1em;
    }

    textarea {
        min-height: 100px;
        resize: vertical;
        padding-right: 0.8rem;
    }

    .select-wrapper {
        position: relative;
        background: #2d3436;
        border: 1px solid #0984e3;
        border-radius: 3px;
    }

    select {
        padding-right: 2rem;
        appearance: none;
        background: transparent;
    }

    .select-wrapper::after {
        content: "▼";
        position: absolute;
        right: 0.8rem;
        top: 50%;
        transform: translateY(-50%);
        color: #00b894;
        font-size: 0.8rem;
        pointer-events: none;
    }

    .task-form-buttons {
        display: flex;
        gap: 1rem;
        justify-content: flex-end;
        margin-top: 0.5rem;
    }

    .btn-cancel {
        border-color: #636e72;
        color: #636e72;
    }

    .form-status {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: 0.7rem;
        color: #636e72;
    }

    .status-indicator {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background: #00b894;
        animation: pulse 2s infinite;
    }

    @keyframes pulse {
        0%, 100% { opacity: 1; }
        50% { opacity: 0.5; }
    }
    @media (max-width: 768px) {
        .quick-add {
            flex-direction: column;
            gap: 0.8rem;
        }

        .input-wrapper {
            width: 100%;
        }

        input, textarea, select {
            padding: 1rem 3.5rem 1rem 0.5rem; /* Increased padding for better touch targets */
            font-size: 1rem; /* Larger font size for mobile */
        }

        .quick-add-buttons {
            width: 100%;
            justify-content: space-between;
        }

        button {
            padding: 0.8rem 1rem; /* Larger touch targets */
            font-size: 0.9rem;
        }

        .prompt {
            font-size: 1rem;
            padding-left: 0.5rem;
        }

        .quick-character-count {
            right: 1rem;
            font-size: 0.8rem;
        }

        /* Adjust full form elements */
        .task-form textarea {
            min-height: 120px; /* Larger textarea */
            font-size: 1rem;
        }

        .select-wrapper select {
            padding: 1rem;
            font-size: 1rem;
        }

        .task-form-buttons {
            flex-direction: row;
            gap: 1rem;
        }

        .task-form-buttons button {
            flex: 1; /* Make buttons fill the width */
        }
    }

    /* For even smaller screens */
    @media (max-width: 480px) {
        .task-creation-container {
            padding: 0.8rem;
        }

        input, textarea, select {
            font-size: 16px; /* Prevents iOS zoom on focus */
        }
    }
</style>