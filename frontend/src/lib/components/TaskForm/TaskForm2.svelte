<script>
    import { browser } from '$app/environment';
    import { api } from '$lib/api.js';
    import { tasks } from '$lib/stores.js';
    import { mdiPlus, mdiChevronDown } from '@mdi/js';
    import { Analytics } from '$lib/analytics';


    let title = '';
    let description = '';
    let status = 'pending';
    let isQuickAdd = true;
    let isSubmitting = false;

    async function handleQuickSubmit() {

        if (!title.trim() || isSubmitting) return;
        
        if (title.length > 100) {
            showError("Task title cannot exceed 100 characters.");
            return;
        }

        Analytics.track('Task Created', {
            type: 'quick',
            title: title
        });

        isSubmitting = true;
        try {
            await api.addTask(title, '', 'pending');
            const updatedTasks = await api.fetchTasks();
            tasks.set(updatedTasks);
            title = '';
        } finally {
            isSubmitting = false;
        }
    }

    async function handleFullSubmit() {

        if (isSubmitting) return;

        if (title.length > 100) {
            showError("Task title cannot exceed 100 characters.");
            return;
        }

        Analytics.track('Task Created', {
            type: 'full',
            title: title,
            description: description,
            status: status
        });

        isSubmitting = true;
        try {
            await api.addTask(title, description, status);
            const updatedTasks = await api.fetchTasks();
            tasks.set(updatedTasks);
            
            title = '';
            description = '';
            status = 'pending';
            isQuickAdd = true;
        } finally {
            isSubmitting = false;
        }
    }
</script>

<div class="task-creation-container" class:submitting={isSubmitting}>
    {#if isQuickAdd}
        <!-- Quick Add Form -->
        <div class="quick-add terminal-box">
            <div class="terminal-header">
                <span class="terminal-dots">
                    <span class="dot"></span>
                    <span class="dot"></span>
                    <span class="dot"></span>
                </span>
                <span class="terminal-title">QUICK_TASK.exe</span>
                <div class="system-status">
                    <span class="status-dot"></span>
                    <span class="status-text">READY</span>
                </div>
            </div>
            
            <div class="terminal-content">
                <div class="input-wrapper">
                    <span class="prompt">>_</span>
                    <input 
                        bind:value={title} 
                        type="text" 
                        placeholder="INITIALIZE_NEW_TASK..." 
                        on:keydown={(e) => e.key === 'Enter' && handleQuickSubmit()}
                        maxlength="100"
                        disabled={isSubmitting}
                    />
                    <span class="character-count" class:warning={title.length > 90}>
                        [{title.length}/100]
                    </span>
                </div>

                <div class="action-buttons">
                    <button 
                        class="cyber-button primary" 
                        on:click={handleQuickSubmit}
                        disabled={!title.trim() || isSubmitting}
                    >
                        <span class="btn-content">
                            <span class="btn-icon">⚡</span>
                            <span class="btn-text">EXECUTE</span>
                        </span>
                        <div class="btn-glow"></div>
                    </button>

                    <button 
                        class="cyber-button secondary" 
                        on:click={() => isQuickAdd = false}
                        disabled={isSubmitting}
                    >
                        <span class="btn-content">
                            <span class="btn-icon">⚙</span>
                            <span class="btn-text">CONFIG</span>
                        </span>
                        <div class="btn-glow"></div>
                    </button>
                </div>
            </div>
        </div>
        {:else}
        <form class="task-form terminal-box" on:submit|preventDefault={handleFullSubmit}>
            <div class="terminal-header">
                <span class="terminal-dots">
                    <span class="dot"></span>
                    <span class="dot"></span>
                    <span class="dot"></span>
                </span>
                <span class="terminal-title">TASK_CONFIGURATION.exe</span>
                <div class="system-status">
                    <span class="status-dot"></span>
                    <span class="status-text">{isSubmitting ? 'PROCESSING' : 'READY'}</span>
                </div>
            </div>
            
            <div class="form-grid">
                <div class="input-group">
                    <div class="input-label">[TASK_IDENTIFIER]</div>
                    <div class="input-wrapper">
                        <span class="prompt">>_</span>
                        <input 
                            bind:value={title} 
                            type="text" 
                            placeholder="Enter task name..." 
                            required
                            maxlength="100"
                            disabled={isSubmitting}
                        />
                        <span class="character-count" class:warning={title.length > 90}>
                            [{title.length}/100]
                        </span>
                    </div>
                </div>
    
                <div class="input-group">
                    <div class="input-label">[TASK_PARAMETERS]</div>
                    <div class="input-wrapper">
                        <span class="prompt">>_</span>
                        <textarea 
                            bind:value={description} 
                            placeholder="Enter task description..."
                            disabled={isSubmitting}
                        ></textarea>
                    </div>
                </div>
    
                <div class="input-group">
                    <div class="input-label">[INITIAL_STATUS]</div>
                    <div class="select-wrapper">
                        <select 
                            bind:value={status} 
                            required
                            disabled={isSubmitting}
                        >
                            <option value="pending">PENDING::0x01</option>
                            <option value="in_progress">IN_PROGRESS::0x02</option>
                            <option value="completed">COMPLETED::0x03</option>
                        </select>
                        <div class="select-arrow">
                            <svg viewBox="0 0 24 24" width="16" height="16">
                                <path 
                                    d="M7 10l5 5 5-5" 
                                    fill="none" 
                                    stroke="currentColor" 
                                    stroke-width="2"
                                    stroke-linecap="round"
                                />
                            </svg>
                        </div>
                    </div>
                </div>
            </div>
    
            <div class="form-actions">
                <button 
                    type="button" 
                    class="cyber-button secondary" 
                    on:click={() => isQuickAdd = true}
                    disabled={isSubmitting}
                >
                    <span class="btn-content">
                        <span class="btn-icon">×</span>
                        <span class="btn-text">ABORT</span>
                    </span>
                    <div class="btn-glow"></div>
                </button>
                
                <button 
                    type="submit" 
                    class="cyber-button primary"
                    disabled={isSubmitting || !title.trim()}
                >
                    <span class="btn-content">
                        <span class="btn-icon">⚡</span>
                        <span class="btn-text">{isSubmitting ? 'PROCESSING...' : 'EXECUTE'}</span>
                    </span>
                    <div class="btn-glow"></div>
                </button>
            </div>
    
            <div class="form-status">
                <div class="status-line">
                    <span class="status-indicator"></span>
                    <span class="status-text">SYSTEM::{isSubmitting ? 'PROCESSING' : 'READY'}</span>
                </div>
                <div class="hex-code">0x{Math.floor(Math.random() * 0xFFFF).toString(16).padStart(4, '0')}</div>
            </div>
        </form>
    {/if}
</div>

<style>
    /* Base Styles */
    .task-creation-container {
        font-family: 'JetBrains Mono', monospace;
        margin-bottom: 2rem;
        position: relative;
    }

    .terminal-box {
        background: #1c1c1c;
        border: 1px solid #0984e3;
        border-radius: 4px;
        overflow: hidden;
        position: relative;
    }

    .terminal-box::before {
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
        background: #2d3436;
        padding: 0.5rem;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        border-bottom: 1px solid rgba(9, 132, 227, 0.2);
    }

    .terminal-dots {
        display: flex;
        gap: 4px;
    }

    .dot {
        width: 6px;
        height: 6px;
        border-radius: 50%;
        background: #636e72;
    }

    .terminal-title {
        color: #00b894;
        font-size: 0.7rem;
        letter-spacing: 0.1em;
    }

    .system-status {
        margin-left: auto;
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }

    .status-dot {
        width: 6px;
        height: 6px;
        border-radius: 50%;
        background: #00b894;
        animation: pulse 2s infinite;
    }

    .status-text {
        color: #00b894;
        font-size: 0.6rem;
        letter-spacing: 0.1em;
    }

    .submitting .status-dot {
        background: #0984e3;
    }

    .submitting .status-text {
        color: #0984e3;
    }

    @keyframes pulse {
        0%, 100% {
            transform: scale(1);
            opacity: 1;
        }
        50% {
            transform: scale(1.2);
            opacity: 0.5;
        }
    }

    /* Quick Add Form Styles */
    .terminal-content {
        padding: 1rem;
    }

    .input-wrapper {
        position: relative;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        background: #2d3436;
        border: 1px solid #0984e3;
        border-radius: 3px;
        padding: 0 0.5rem;
        margin-bottom: 1rem;
        transition: all 0.3s ease;
    }

    .input-wrapper:focus-within {
        border-color: #00b894;
        box-shadow: 0 0 8px rgba(0, 184, 148, 0.2);
    }

    .prompt {
        color: #00b894;
        font-size: 0.9rem;
        user-select: none;
    }

    input {
        flex: 1;
        background: transparent;
        border: none;
        color: #fff;
        padding: 0.8rem 3.5rem 0.8rem 0.5rem;
        font-family: inherit;
        font-size: 0.9rem;
    }

    input:focus {
        outline: none;
    }

    input:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .character-count {
        position: absolute;
        right: 0.8rem;
        color: #636e72;
        font-size: 0.7rem;
        pointer-events: none;
    }

    .character-count.warning {
        color: #ff6b6b;
        animation: blink 1s infinite;
    }

    .action-buttons {
        display: flex;
        gap: 0.8rem;
    }

    .cyber-button {
        flex: 1;
        position: relative;
        background: transparent;
        border: 1px solid currentColor;
        padding: 0.8rem;
        font-family: inherit;
        font-size: 0.8rem;
        cursor: pointer;
        overflow: hidden;
        transition: all 0.3s ease;
        border-radius: 3px;
    }

    .btn-content {
        position: relative;
        z-index: 2;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
    }

    .btn-glow {
        position: absolute;
        top: 0;
        left: -100%;
        width: 100%;
        height: 100%;
        background: linear-gradient(
            90deg,
            transparent,
            rgba(255, 255, 255, 0.1),
            transparent
        );
        transition: 0.5s;
    }

    .cyber-button:hover .btn-glow {
        left: 100%;
    }

    .cyber-button.primary {
        color: #00b894;
    }

    .cyber-button.primary:hover {
        background: rgba(0, 184, 148, 0.1);
        box-shadow: 0 0 10px rgba(0, 184, 148, 0.3);
    }

    .cyber-button.secondary {
        color: #0984e3;
    }

    .cyber-button.secondary:hover {
        background: rgba(9, 132, 227, 0.1);
        box-shadow: 0 0 10px rgba(9, 132, 227, 0.3);
    }

    .cyber-button:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .cyber-button:disabled .btn-glow {
        display: none;
    }

    /* Scanline effect */
    .terminal-box::after {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        height: 1px;
        background: linear-gradient(
            90deg,
            transparent,
            rgba(0, 184, 148, 0.2),
            transparent
        );
        animation: scan 2s linear infinite;
        pointer-events: none;
    }

    @keyframes scan {
        from { transform: translateY(-100%); }
        to { transform: translateY(100%); }
    }

    @keyframes blink {
        0%, 100% { opacity: 1; }
        50% { opacity: 0.5; }
    }

    /* Loading state */
    .submitting .input-wrapper {
        opacity: 0.7;
    }

    .submitting .cyber-button:not(:disabled) {
        position: relative;
        overflow: hidden;
    }

    .submitting .cyber-button:not(:disabled)::after {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: linear-gradient(
            90deg,
            transparent,
            rgba(255, 255, 255, 0.1),
            transparent
        );
        animation: loading 1s infinite;
    }

    @keyframes loading {
        from { transform: translateX(-100%); }
        to { transform: translateX(100%); }
    }

    /* Full Form Styles */
    .form-grid {
        display: flex;
        flex-direction: column;
        gap: 1.5rem;
        padding: 1.5rem;
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
        flex: 1;
        background: transparent;
        border: none;
        color: #fff;
        padding: 0.8rem;
        font-family: inherit;
        font-size: 0.9rem;
        resize: vertical;
        min-height: 100px;
    }

    textarea:focus {
        outline: none;
    }

    .select-wrapper {
        position: relative;
        background: #2d3436;
        border: 1px solid #0984e3;
        border-radius: 3px;
        transition: all 0.3s ease;
    }

    .select-wrapper:focus-within {
        border-color: #00b894;
        box-shadow: 0 0 8px rgba(0, 184, 148, 0.2);
    }

    select {
        width: 100%;
        background: transparent;
        border: none;
        color: #fff;
        padding: 0.8rem;
        font-family: inherit;
        font-size: 0.9rem;
        appearance: none;
        cursor: pointer;
    }

    select:focus {
        outline: none;
    }

    .select-arrow {
        position: absolute;
        right: 0.8rem;
        top: 50%;
        transform: translateY(-50%);
        color: #0984e3;
        pointer-events: none;
        transition: transform 0.3s ease;
    }

    .select-wrapper:focus-within .select-arrow {
        transform: translateY(-50%) rotate(180deg);
        color: #00b894;
    }

    .form-actions {
        display: flex;
        gap: 1rem;
        padding: 0 1.5rem 1.5rem;
    }

    .form-status {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 0.8rem 1.5rem;
        background: rgba(45, 52, 54, 0.5);
        border-top: 1px solid rgba(9, 132, 227, 0.2);
    }

    .status-line {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }

    .status-indicator {
        width: 6px;
        height: 6px;
        border-radius: 50%;
        background: #00b894;
        animation: pulse 2s infinite;
    }

    .hex-code {
        color: #636e72;
        font-size: 0.7rem;
        letter-spacing: 0.1em;
    }

    /* Processing state styles */
    .submitting .status-indicator {
        background: #0984e3;
    }

    .submitting textarea,
    .submitting select {
        opacity: 0.7;
    }

    /* Circuit board pattern for form */
    .form-grid::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background-image: 
            linear-gradient(rgba(0, 184, 148, 0.05) 1px, transparent 1px),
            linear-gradient(90deg, rgba(0, 184, 148, 0.05) 1px, transparent 1px);
        background-size: 20px 20px;
        pointer-events: none;
        opacity: 0.5;
    }
</style>