<script>
    import { mdiCheckCircle, mdiClockOutline, mdiDragVertical, mdiChartBar, mdiStar } from '@mdi/js';
    import { onMount } from 'svelte';

    let isVisible = false;

    onMount(() => {
        isVisible = true;
    });
</script>

<div class="landing-container">
    <header class="hero">
        <div class="hero-content">
            <h1 class:animate-slide-up={isVisible}>ActionHub</h1>
            <p class="tagline" class:animate-slide-up-delay={isVisible}>
                Streamline your productivity with our intuitive task management solution
            </p>
            <div class="cta-buttons" class:animate-fade-in={isVisible}>
                <a href="/register" class="btn btn-primary">Get Started Free</a>
                <a href="/login" class="btn btn-secondary">Sign In</a>
            </div>
        </div>
    </header>

    <section class="stats">
        <div class="stat-item">
            <div class="stat-number">10k+</div>
            <div class="stat-label">Active Users</div>
        </div>
        <div class="stat-item">
            <div class="stat-number">1M+</div>
            <div class="stat-label">Tasks Completed</div>
        </div>
        <div class="stat-item">
            <div class="stat-number">99%</div>
            <div class="stat-label">Satisfaction Rate</div>
        </div>
    </section>

    <section class="features">
        <h2>Why Choose ActionHub?</h2>
        <div class="feature-grid">
            <div class="feature-card">
                <div class="feature-icon">
                    <svg viewBox="0 0 24 24" width="24" height="24">
                        <path fill="currentColor" d={mdiCheckCircle} />
                    </svg>
                </div>
                <h3>Simple & Intuitive</h3>
                <p>Quick task creation and easy organization with our user-friendly interface</p>
            </div>

            <div class="feature-card">
                <div class="feature-icon">
                    <svg viewBox="0 0 24 24" width="24" height="24">
                        <path fill="currentColor" d={mdiDragVertical} />
                    </svg>
                </div>
                <h3>Drag & Drop</h3>
                <p>Effortlessly prioritize tasks with drag-and-drop functionality</p>
            </div>

            <div class="feature-card">
                <div class="feature-icon">
                    <svg viewBox="0 0 24 24" width="24" height="24">
                        <path fill="currentColor" d={mdiClockOutline} />
                    </svg>
                </div>
                <h3>Real-time Updates</h3>
                <p>Stay on top of your tasks with instant status updates</p>
            </div>

            <div class="feature-card">
                <div class="feature-icon">
                    <svg viewBox="0 0 24 24" width="24" height="24">
                        <path fill="currentColor" d={mdiChartBar} />
                    </svg>
                </div>
                <h3>Progress Tracking</h3>
                <p>Monitor your productivity with detailed statistics</p>
            </div>
        </div>
    </section>

    <section class="testimonials">
        <h2>What Our Users Say</h2>
        <div class="testimonial-grid">
            <div class="testimonial-card">
                <div class="testimonial-rating">
                    {#each Array(5) as _}
                        <svg viewBox="0 0 24 24" width="20" height="20">
                            <path fill="currentColor" d={mdiStar} />
                        </svg>
                    {/each}
                </div>
                <p class="testimonial-text">"ActionHub has transformed how I manage my daily tasks. The interface is intuitive and the features are exactly what I need."</p>
                <div class="testimonial-author">
                    <div class="author-avatar">JD</div>
                    <div class="author-info">
                        <strong>John Doe</strong>
                        <span>Product Manager</span>
                    </div>
                </div>
            </div>
            
            <div class="testimonial-card">
                <div class="testimonial-rating">
                    {#each Array(5) as _}
                        <svg viewBox="0 0 24 24" width="20" height="20">
                            <path fill="currentColor" d={mdiStar} />
                        </svg>
                    {/each}
                </div>
                <p class="testimonial-text">"The best task management tool I've ever used. Simple yet powerful!"</p>
                <div class="testimonial-author">
                    <div class="author-avatar">AS</div>
                    <div class="author-info">
                        <strong>Alice Smith</strong>
                        <span>Freelancer</span>
                    </div>
                </div>
            </div>
        </div>
    </section>

    <section class="cta">
        <h2>Ready to Get Started?</h2>
        <p>Join thousands of users who have improved their productivity</p>
        <a href="/register" class="btn btn-primary">Create Free Account</a>
    </section>
</div>

<style>
    .landing-container {
        max-width: 1200px;
        margin: 0 auto;
        padding: 20px;
        font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
    }

    .hero {
        text-align: center;
        padding: 80px 20px;
        background: linear-gradient(135deg, #f6f9fc 0%, #efe9ef 100%);
        border-radius: 16px;
        margin-bottom: 40px;
        position: relative;
        overflow: hidden;
    }

    .hero h1 {
        font-size: 4rem;
        margin-bottom: 1rem;
        color: #2c3e50;
        opacity: 0;
    }

    .tagline {
        font-size: 1.5rem;
        color: #6c757d;
        margin-bottom: 2rem;
        opacity: 0;
    }

    .cta-buttons {
        display: flex;
        gap: 1rem;
        justify-content: center;
        opacity: 0;
    }

    .animate-slide-up {
        animation: slideUp 0.6s ease forwards;
    }

    .animate-slide-up-delay {
        animation: slideUp 0.6s ease forwards;
        animation-delay: 0.2s;
    }

    .animate-fade-in {
        animation: fadeIn 0.6s ease forwards;
        animation-delay: 0.4s;
    }

    @keyframes slideUp {
        from { transform: translateY(20px); opacity: 0; }
        to { transform: translateY(0); opacity: 1; }
    }

    @keyframes fadeIn {
        from { opacity: 0; }
        to { opacity: 1; }
    }

    .btn {
        padding: 14px 28px;
        border-radius: 8px;
        text-decoration: none;
        font-weight: 600;
        transition: all 0.2s;
        font-size: 1.1rem;
    }

    .btn-primary {
        background: linear-gradient(45deg, #4CAF50, #45a049);
        color: white;
        border: none;
        box-shadow: 0 4px 15px rgba(76, 175, 80, 0.3);
    }

    .btn-primary:hover {
        transform: translateY(-2px);
        box-shadow: 0 6px 20px rgba(76, 175, 80, 0.4);
    }

    .btn-secondary {
        background-color: white;
        color: #4CAF50;
        border: 2px solid #4CAF50;
    }

    .btn-secondary:hover {
        background-color: #f8f9fa;
        transform: translateY(-2px);
    }

    .stats {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
        gap: 2rem;
        padding: 4rem 2rem;
        text-align: center;
        background: #f8f9fa;
        border-radius: 16px;
        margin-bottom: 40px;
    }

    .stat-number {
        font-size: 2.5rem;
        font-weight: bold;
        color: #4CAF50;
        margin-bottom: 0.5rem;
    }

    .stat-label {
        font-size: 1.1rem;
        color: #6c757d;
    }

    .features {
        padding: 60px 0;
    }

    .features h2 {
        text-align: center;
        margin-bottom: 40px;
        color: #2c3e50;
        font-size: 2.5rem;
    }

    .feature-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
        gap: 30px;
        padding: 20px;
    }

    .feature-card {
        padding: 40px;
        background: white;
        border-radius: 16px;
        box-shadow: 0 8px 32px rgba(31, 38, 135, 0.15);
        text-align: center;
        transition: all 0.3s ease;
        border: 1px solid rgba(255, 255, 255, 0.2);
    }

    .feature-card:hover {
        transform: translateY(-10px);
        box-shadow: 0 12px 40px rgba(31, 38, 135, 0.2);
    }

    .feature-icon {
        margin-bottom: 20px;
    }

    .feature-icon svg {
        width: 60px;
        height: 60px;
        color: #4CAF50;
    }

    .feature-card h3 {
        margin-bottom: 15px;
        color: #2c3e50;
        font-size: 1.5rem;
    }

    .feature-card p {
        color: #6c757d;
        line-height: 1.6;
        font-size: 1.1rem;
    }

    .testimonials {
        padding: 60px 0;
    }

    .testimonials h2 {
        text-align: center;
        margin-bottom: 40px;
        color: #2c3e50;
        font-size: 2.5rem;
    }

    .testimonial-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
        gap: 30px;
        padding: 20px;
    }

    .testimonial-card {
        padding: 30px;
        background: white;
        border-radius: 16px;
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
    }

    .testimonial-rating {
        color: #ffd700;
        margin-bottom: 20px;
        display: flex;
        gap: 4px;
    }

    .testimonial-text {
        color: #2c3e50;
        font-size: 1.1rem;
        line-height: 1.6;
        margin-bottom: 20px;
        font-style: italic;
    }

    .testimonial-author {
        display: flex;
        align-items: center;
        gap: 15px;
    }

    .author-avatar {
        width: 50px;
        height: 50px;
        background: #4CAF50;
        color: white;
        border-radius: 25px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-weight: bold;
    }

    .author-info {
        display: flex;
        flex-direction: column;
    }

    .author-info strong {
        color: #2c3e50;
    }

    .author-info span {
        color: #6c757d;
        font-size: 0.9rem;
    }

    .cta {
        text-align: center;
        padding: 80px 20px;
        background: linear-gradient(135deg, #f6f9fc 0%, #e9ecef 100%);
        border-radius: 16px;
        margin-top: 40px;
    }

    .cta h2 {
        margin-bottom: 1rem;
        color: #2c3e50;
        font-size: 2.5rem;
    }

    .cta p {
        margin-bottom: 2rem;
        color: #6c757d;
        font-size: 1.2rem;
    }

    @media (max-width: 768px) {
        .hero h1 {
            font-size: 2.5rem;
        }

        .tagline {
            font-size: 1.2rem;
        }

        .feature-grid {
            grid-template-columns: 1fr;
        }

        .cta-buttons {
            flex-direction: column;
            align-items: center;
        }

        .btn {
            width: 100%;
            max-width: 300px;
            text-align: center;
        }

        .stats {
            grid-template-columns: 1fr;
            gap: 2rem;
        }

        .testimonial-grid {
            grid-template-columns: 1fr;
        }
    }
</style>