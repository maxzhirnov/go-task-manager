<script>
    import { mdiCheckCircle, mdiClockOutline, mdiDragVertical, mdiChartBar, mdiStar } from '@mdi/js';
    import { onMount } from 'svelte';

    let isVisible = false;

    onMount(() => {
        isVisible = true;
    });
</script>

<div class="landing-container">
    <div class="cyber-grid"></div>
    
    <header class="hero">
        <div class="glitch-container">
            <h1 class="glitch" class:animate-slide-up={isVisible}>
                <span aria-hidden="true">ActionHub</span>
                ActionHub
                <span aria-hidden="true">ActionHub</span>
            </h1>
        </div>
        
        <div class="terminal-box hero-terminal">
            <div class="terminal-header">
                <span class="terminal-dots">
                    <span class="dot"></span>
                    <span class="dot"></span>
                    <span class="dot"></span>
                </span>
                <span class="terminal-title">SYSTEM_OVERVIEW.exe</span>
            </div>
            
            <div class="terminal-content">
                <div class="typing-text terminal-content-text" class:animate-slide-up-delay={isVisible}>
                    <span class="prompt">>_</span> INITIALIZING PRODUCTIVITY SYSTEM
                    <br/>
                    <span class="prompt">>_</span> LOADING CORE MODULES...
                    <br/>
                    <span class="prompt">>_</span> SYSTEM READY
                </div>
                
                <div class="cta-buttons" class:animate-fade-in={isVisible}>
                    <a href="/register" class="cyber-button">
                        <span class="btn-icon">⚡</span>
                        <span class="btn-text">INITIALIZE_SESSION</span>
                    </a>
                    <a href="/login" class="cyber-button secondary">
                        <span class="btn-icon">🔒</span>
                        <span class="btn-text">ACCESS_PORTAL</span>
                    </a>
                </div>
            </div>
        </div>
    </header>

    <section class="stats">
        <div class="terminal-box">
            <div class="terminal-header">
                <span class="terminal-title">SYSTEM_METRICS.dat</span>
            </div>
            <div class="stats-grid">
                <div class="stat-item">
                    <div class="stat-number">10k+</div>
                    <div class="stat-label">ACTIVE_USERS</div>
                    <div class="stat-pulse"></div>
                </div>
                <div class="stat-item">
                    <div class="stat-number">1M+</div>
                    <div class="stat-label">TASKS_COMPLETED</div>
                    <div class="stat-pulse"></div>
                </div>
                <div class="stat-item">
                    <div class="stat-number">99%</div>
                    <div class="stat-label">SATISFACTION_RATE</div>
                    <div class="stat-pulse"></div>
                </div>
            </div>
        </div>
    </section>
        <!-- Features Section -->
        <section class="features">
            <div class="terminal-box">
                <div class="terminal-header">
                    <span class="terminal-dots">
                        <span class="dot"></span>
                        <span class="dot"></span>
                        <span class="dot"></span>
                    </span>
                    <span class="terminal-title">CORE_FEATURES.sys</span>
                </div>
                
                <div class="features-grid">
                    <div class="feature-card">
                        <div class="feature-icon">
                            <svg viewBox="0 0 24 24" width="24" height="24">
                                <path fill="currentColor" d={mdiCheckCircle} />
                            </svg>
                        </div>
                        <div class="feature-content">
                            <h3>[INTUITIVE_INTERFACE]</h3>
                            <div class="feature-description">
                                <span class="prompt">>_</span>
                                Quick task creation and streamlined organization protocols
                            </div>
                        </div>
                    </div>
    
                    <div class="feature-card">
                        <div class="feature-icon">
                            <svg viewBox="0 0 24 24" width="24" height="24">
                                <path fill="currentColor" d={mdiDragVertical} />
                            </svg>
                        </div>
                        <div class="feature-content">
                            <h3>[DYNAMIC_CONTROL]</h3>
                            <div class="feature-description">
                                <span class="prompt">>_</span>
                                Advanced drag-drop task prioritization system
                            </div>
                        </div>
                    </div>
    
                    <div class="feature-card">
                        <div class="feature-icon">
                            <svg viewBox="0 0 24 24" width="24" height="24">
                                <path fill="currentColor" d={mdiClockOutline} />
                            </svg>
                        </div>
                        <div class="feature-content">
                            <h3>[REAL_TIME_SYNC]</h3>
                            <div class="feature-description">
                                <span class="prompt">>_</span>
                                Instant status updates and synchronization
                            </div>
                        </div>
                    </div>
    
                    <div class="feature-card">
                        <div class="feature-icon">
                            <svg viewBox="0 0 24 24" width="24" height="24">
                                <path fill="currentColor" d={mdiChartBar} />
                            </svg>
                        </div>
                        <div class="feature-content">
                            <h3>[METRICS_ENGINE]</h3>
                            <div class="feature-description">
                                <span class="prompt">>_</span>
                                Advanced productivity analytics and tracking
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>
    
        <!-- Testimonials Section -->
        <section class="testimonials">
            <div class="terminal-box">
                <div class="terminal-header">
                    <span class="terminal-dots">
                        <span class="dot"></span>
                        <span class="dot"></span>
                        <span class="dot"></span>
                    </span>
                    <span class="terminal-title">USER_FEEDBACK.log</span>
                </div>
                
                <div class="testimonials-grid">
                    <div class="testimonial-card">
                        <div class="testimonial-header">
                            <span class="user-id">USER::JD_001</span>
                            <div class="rating">
                                {#each Array(5) as _}
                                    <span class="star">★</span>
                                {/each}
                            </div>
                        </div>
                        <div class="testimonial-content">
                            <span class="prompt">>_</span>
                            "ActionHub has transformed my task management protocols. Interface optimization at its finest."
                        </div>
                        <div class="user-meta">
                            <div class="user-avatar">JD</div>
                            <div class="user-info">
                                <span class="user-name">John.Doe</span>
                                <span class="user-role">Product_Manager.exe</span>
                            </div>
                        </div>
                    </div>
                    
                    <div class="testimonial-card">
                        <div class="testimonial-header">
                            <span class="user-id">USER::AS_002</span>
                            <div class="rating">
                                {#each Array(5) as _}
                                    <span class="star">★</span>
                                {/each}
                            </div>
                        </div>
                        <div class="testimonial-content">
                            <span class="prompt">>_</span>
                            "Maximum efficiency achieved. Best task management system in the matrix."
                        </div>
                        <div class="user-meta">
                            <div class="user-avatar">AS</div>
                            <div class="user-info">
                                <span class="user-name">Alice.Smith</span>
                                <span class="user-role">Freelance.sys</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>
    
        <!-- Final CTA Section -->
        <section class="final-cta">
            <div class="terminal-box">
                <div class="terminal-header">
                    <span class="terminal-title">INITIALIZE_SESSION.exe</span>
                </div>
                <div class="cta-content">
                    <div class="typing-text cta-content-text">
                        <span class="prompt">>_</span> READY TO UPGRADE YOUR PRODUCTIVITY?
                        <br/>
                        <span class="prompt">>_</span> JOIN THE NETWORK
                    </div>
                    <a href="/register" class="cyber-button">
                        <span class="btn-icon">⚡</span>
                        <span class="btn-text">CREATE_NEW_INSTANCE</span>
                    </a>
                </div>
            </div>
        </section>
    </div>
    
    <style>

        section {
            padding: 0 1rem;
        }

        .landing-container {
            background: #0f1215;
            color: #fff;
            font-family: "JetBrains Mono", monospace;
            position: relative;
            overflow: hidden;
        }
    
        .cyber-grid {
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background-image: 
                linear-gradient(rgba(9, 132, 227, 0.1) 1px, transparent 1px),
                linear-gradient(90deg, rgba(9, 132, 227, 0.1) 1px, transparent 1px);
            background-size: 30px 30px;
            z-index: 1;
            pointer-events: none;
        }
    
        .hero {
            min-height: 80vh;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            padding: 2rem;
            position: relative;
            z-index: 2;
        }
    
        .glitch {
            font-size: 4rem;
            font-weight: bold;
            text-transform: uppercase;
            position: relative;
            text-shadow: 0.05em 0 0 #00fffc, -0.03em -0.04em 0 #fc00ff,
                         0.025em 0.04em 0 #fffc00;
            animation: glitch 725ms infinite;
        }
    
        .glitch span {
            position: absolute;
            top: 0;
            left: 0;
        }
    
        .glitch span:first-child {
            animation: glitch 500ms infinite;
            clip-path: polygon(0 0, 100% 0, 100% 35%, 0 35%);
            transform: translate(-0.04em, -0.03em);
            opacity: 0.75;
        }
    
        .glitch span:last-child {
            animation: glitch 375ms infinite;
            clip-path: polygon(0 65%, 100% 65%, 100% 100%, 0 100%);
            transform: translate(0.04em, 0.03em);
            opacity: 0.75;
        }
    
        .terminal-box {
            background: #1c1c1c;
            border: 1px solid #0984e3;
            border-radius: 4px;
            overflow: hidden;
            position: relative;
            width: 100%;
            max-width: 800px;
            margin: 2rem auto;
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

    .terminal-content {
        padding: 2rem;
        position: relative;
    }

    .terminal-content-text {
        margin-bottom: 1rem;
    }

    .typing-text {
        color: #00b894;
        line-height: 1.6;
    }

    .prompt {
        color: #0984e3;
        margin-right: 0.5rem;
    }

    /* Features Styling */
    .features-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
        gap: 2rem;
        padding: 2rem;
    }

    .feature-card {
        background: #2d3436;
        border: 1px solid #0984e3;
        border-radius: 4px;
        padding: 1.5rem;
        transition: all 0.3s ease;
        position: relative;
        overflow: hidden;
    }

    .feature-card::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: linear-gradient(45deg, transparent 0%, rgba(9, 132, 227, 0.1) 100%);
        opacity: 0;
        transition: opacity 0.3s ease;
    }

    .feature-card:hover::before {
        opacity: 1;
    }

    .feature-icon {
        color: #00b894;
        margin-bottom: 1rem;
    }

    .feature-content h3 {
        color: #0984e3;
        font-size: 0.9rem;
        margin-bottom: 0.5rem;
        letter-spacing: 0.1em;
    }

    .feature-description {
        color: #fff;
        font-size: 0.8rem;
        line-height: 1.6;
    }

    /* Stats Styling */
    .stats-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
        gap: 2rem;
        padding: 2rem;
    }

    .stat-item {
        text-align: center;
        position: relative;
    }

    .stat-number {
        color: #00b894;
        font-size: 2.5rem;
        font-weight: bold;
        margin-bottom: 0.5rem;
        text-shadow: 0 0 10px rgba(0, 184, 148, 0.5);
    }

    .stat-label {
        color: #0984e3;
        font-size: 0.8rem;
        letter-spacing: 0.1em;
    }

    .stat-pulse {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        width: 100px;
        height: 100px;
        border: 2px solid #00b894;
        border-radius: 50%;
        opacity: 0;
        animation: pulse 2s infinite;
    }

    /* Testimonials Styling */
    .testimonials-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
        gap: 2rem;
        padding: 2rem;
    }

    .testimonial-card {
        background: #2d3436;
        border: 1px solid #0984e3;
        border-radius: 4px;
        padding: 1.5rem;
        position: relative;
    }

    .testimonial-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1rem;
    }

    .user-id {
        color: #0984e3;
        font-size: 0.7rem;
        letter-spacing: 0.1em;
    }

    .rating {
        color: #00b894;
    }

    .star {
        margin-left: 2px;
    }

    .testimonial-content {
        color: #fff;
        font-size: 0.9rem;
        line-height: 1.6;
        margin-bottom: 1rem;
    }

    .user-meta {
        display: flex;
        align-items: center;
        gap: 1rem;
    }

    .user-avatar {
        width: 40px;
        height: 40px;
        background: #0984e3;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #fff;
        font-size: 0.8rem;
    }

    .user-info {
        display: flex;
        flex-direction: column;
    }

    .user-name {
        color: #00b894;
        font-size: 0.8rem;
    }

    .user-role {
        color: #636e72;
        font-size: 0.7rem;
    }

    /* CTA Styling */
    .cyber-button {
        background: transparent;
        border: 1px solid #00b894;
        color: #00b894;
        padding: 1rem 2rem;
        border-radius: 4px;
        text-decoration: none;
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        font-family: inherit;
        font-size: 0.9rem;
        transition: all 0.3s ease;
        position: relative;
        overflow: hidden;
    }

    .cyber-button::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: linear-gradient(45deg, transparent, rgba(0, 184, 148, 0.2), transparent);
        transform: translateX(-100%);
        transition: transform 0.5s ease;
    }

    .cyber-button:hover::before {
        transform: translateX(100%);
    }

    .cyber-button.secondary {
        border-color: #0984e3;
        color: #0984e3;
    }

    .cta-content {
        padding: 2rem;
    }

    .cta-content-text {
        margin-bottom: 1rem;
    }

    /* Animations */
    @keyframes glitch {
        0% { text-shadow: 0.05em 0 0 #00fffc, -0.03em -0.04em 0 #fc00ff,
                         0.025em 0.04em 0 #fffc00; }
        15% { text-shadow: 0.05em 0 0 #00fffc, -0.03em -0.04em 0 #fc00ff,
                          0.025em 0.04em 0 #fffc00; }
        16% { text-shadow: -0.05em -0.025em 0 #00fffc, 0.025em 0.035em 0 #fc00ff,
                          -0.05em -0.05em 0 #fffc00; }
        49% { text-shadow: -0.05em -0.025em 0 #00fffc, 0.025em 0.035em 0 #fc00ff,
                          -0.05em -0.05em 0 #fffc00; }
        50% { text-shadow: 0.05em 0.035em 0 #00fffc, 0.03em 0 0 #fc00ff,
                          0 -0.04em 0 #fffc00; }
        99% { text-shadow: 0.05em 0.035em 0 #00fffc, 0.03em 0 0 #fc00ff,
                          0 -0.04em 0 #fffc00; }
        100% { text-shadow: -0.05em 0 0 #00fffc, -0.025em -0.04em 0 #fc00ff,
                           -0.04em -0.025em 0 #fffc00; }
    }

    @keyframes pulse {
        0% { transform: translate(-50%, -50%) scale(1); opacity: 0; }
        50% { opacity: 0.3; }
        100% { transform: translate(-50%, -50%) scale(1.5); opacity: 0; }
    }

    /* Responsive Design */
    @media (max-width: 768px) {
        .hero {
            padding: 1rem;
        }

        .glitch {
            font-size: 2.5rem;
        }

        .features-grid,
        .testimonials-grid {
            grid-template-columns: 1fr;
        }

        .terminal-content {
            padding: 1rem;
        }

        .cyber-button {
            width: 80%;
            margin-bottom: .5rem;
            justify-content: center;
        }
    }

    /* Animation Classes */
    .animate-slide-up {
        animation: slideUp 0.8s ease forwards;
    }

    .animate-slide-up-delay {
        animation: slideUp 0.8s ease forwards 0.3s;
        opacity: 0;
    }

    .animate-fade-in {
        animation: fadeIn 1s ease forwards 0.6s;
        opacity: 0;
    }

    @keyframes slideUp {
        from {
            transform: translateY(30px);
            opacity: 0;
        }
        to {
            transform: translateY(0);
            opacity: 1;
        }
    }

    @keyframes fadeIn {
        from { opacity: 0; }
        to { opacity: 1; }
    }
    </style>