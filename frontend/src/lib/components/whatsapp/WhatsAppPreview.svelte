<script>
    import { t } from 'svelte-i18n';
    import { fetchDisclaimer } from '$lib/utils/api.js';

    // Constants for excerpt formatting
    const MAX_EXCERPT_LENGTH = 100;

    // Helper to truncate string at word boundary
    function truncateString(str, maxLength) {
        if (!str || str.length <= maxLength) return str;
        if (maxLength < 3) return str.substring(0, maxLength) + '...';

        const truncated = str.substring(0, maxLength).replace(/\s+$/, '');
        const lastSpace = truncated.lastIndexOf(' ');
        if (lastSpace > maxLength / 2) {
            return truncated.substring(0, lastSpace) + '...';
        }
        return truncated + '...';
    }

    // Generate article URL from slug
    function getArticleURL(slug) {
        return slug ? `https://prima.app/artikel/${slug}` : '';
    }

    // Generate YouTube URL from video ID
    function getYouTubeURL(youtubeId) {
        return youtubeId ? `https://youtube.com/watch?v=${youtubeId}` : '';
    }

    let {
        message = '',
        patientName = '',
        reminderTitle = '',
        timestamp = $bindable(),
        showDisclaimer = true,
        isScheduled = false,
        scheduledTime = null,
        attachments = []
    } = $props();

    // Disclaimer config from backend (with fallback to i18n)
    let disclaimerConfig = $state({ text: '', enabled: true, loaded: false });

    // Reactive timestamp that updates every minute for real-time preview
    let currentTime = $state(new Date());

    $effect(() => {
        // Update current time every minute for real-time preview
        const interval = setInterval(() => {
            currentTime = new Date();
        }, 60000);
        return () => clearInterval(interval);
    });

    // Fetch disclaimer from backend with proper cleanup to avoid memory leaks
    $effect(() => {
        let aborted = false;
        const controller = new AbortController();

        fetchDisclaimer({ signal: controller.signal })
            .then(config => {
                if (!aborted) {
                    disclaimerConfig = { ...config, loaded: true };
                }
            })
            .catch(err => {
                if (!aborted && err.name !== 'AbortError') {
                    // Fallback to i18n if API fails
                    disclaimerConfig = { text: '', enabled: true, loaded: true };
                }
            });

        return () => {
            aborted = true;
            controller.abort();
        };
    });

    // Format timestamp to HH:mm - use currentTime for reactivity
    let formattedTime = $derived.by(() => {
        const date = timestamp || currentTime;
        return date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' });
    });

    // Format scheduled time for display
    let formattedScheduledTime = $derived.by(() => {
        if (!scheduledTime) return '06:00';
        try {
            const date = new Date(scheduledTime);
            return date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', hour12: false });
        } catch {
            return '06:00';
        }
    });

    // Greeting text (i18n)
    let greetingText = $derived($t('reminder.preview.greeting', { values: { name: patientName } }));

    // Disclaimer text - use backend config if loaded, otherwise fallback to i18n
    let disclaimerText = $derived(
        disclaimerConfig.loaded && disclaimerConfig.text
            ? disclaimerConfig.text
            : $t('reminder.preview.disclaimer')
    );

    // Should show disclaimer - check both prop and backend config
    let shouldShowDisclaimer = $derived(showDisclaimer && disclaimerConfig.enabled);

    // Build complete message preview with excerpts - matches backend FormatReminderMessageWithExcerpts
    let messageLines = $derived.by(() => {
        let lines = [];

        // Greeting
        if (patientName) {
            lines.push({ type: 'text', content: greetingText });
            lines.push({ type: 'empty' });
        }

        // Title (bold in WhatsApp)
        if (reminderTitle) {
            lines.push({ type: 'bold', content: reminderTitle });
            lines.push({ type: 'empty' });
        }

        // Custom message
        if (message) {
            lines.push({ type: 'text', content: message });
        }

        // Attached content with excerpts
        if (attachments && attachments.length > 0) {
            lines.push({ type: 'empty' });
            lines.push({ type: 'divider' });
            lines.push({ type: 'text', content: $t('reminder.preview.attachmentsHeader') || 'Konten Edukasi:' });

            attachments.forEach((attachment) => {
                if (attachment.type === 'article') {
                    // Format: 📖 Title\nExcerpt...\n🔗 link
                    lines.push({ type: 'article-title', content: `📖 ${attachment.title || 'Konten'}`, attachment });

                    // Add excerpt if available (truncated to 100 chars)
                    if (attachment.excerpt) {
                        const truncatedExcerpt = truncateString(attachment.excerpt, MAX_EXCERPT_LENGTH);
                        lines.push({ type: 'article-excerpt', content: truncatedExcerpt, attachment });
                    }

                    // Add link
                    const articleURL = attachment.slug ? getArticleURL(attachment.slug) : (attachment.url || '');
                    if (articleURL) {
                        lines.push({ type: 'article-link', content: `🔗 ${articleURL}`, attachment, isLink: true });
                    }
                } else if (attachment.type === 'video') {
                    // Format: 🎬 Title\n🔗 link
                    lines.push({ type: 'video-title', content: `🎬 ${attachment.title || 'Konten'}`, attachment });

                    // Add link
                    const youtubeURL = attachment.youtubeId ? getYouTubeURL(attachment.youtubeId) : (attachment.url || '');
                    if (youtubeURL) {
                        lines.push({ type: 'video-link', content: `🔗 ${youtubeURL}`, attachment, isLink: true });
                    }
                }
                // Add empty line between attachments
                lines.push({ type: 'empty' });
            });
        }

        // Disclaimer
        if (shouldShowDisclaimer) {
            lines.push({ type: 'empty' });
            lines.push({ type: 'divider' });
            lines.push({ type: 'italic', content: disclaimerText });
        }

        return lines;
    });
</script>

<div
    class="bg-slate-100 p-4 border-t border-slate-200"
    role="region"
    aria-label={$t('reminder.preview.label')}
>
    <p class="text-xs text-slate-500 mb-2">{$t('reminder.preview.title')}</p>

    <div class="flex justify-end">
        <div
            class="whatsapp-bubble bg-whatsapp-bubble rounded-lg p-3 max-w-[80%] relative text-sm text-slate-900"
            aria-live="polite"
        >
            <!-- Message content with WhatsApp formatting -->
            <div class="whitespace-pre-wrap break-words">
                {#each messageLines as line}
                    {#if line.type === 'bold'}
                        <strong>{line.content}</strong>
                    {:else if line.type === 'italic'}
                        <em class="text-slate-600 text-xs">{line.content}</em>
                    {:else if line.type === 'divider'}
                        <hr class="my-2 border-slate-300" />
                    {:else if line.type === 'empty'}
                        <br />
                    {:else if line.type === 'article-title'}
                        <div class="text-teal-800 font-medium">{line.content}</div>
                    {:else if line.type === 'article-excerpt'}
                        <div class="text-slate-600 text-sm">{line.content}</div>
                    {:else if line.type === 'article-link'}
                        <div class="text-teal-600 text-sm">
                            <a
                                href={line.content.replace('🔗 ', '')}
                                target="_blank"
                                rel="noopener noreferrer"
                                class="hover:underline focus:outline-none focus:ring-2 focus:ring-teal-500 rounded cursor-pointer"
                                tabindex="0"
                                aria-label="Baca artikel lengkap: {line.attachment?.title || 'Artikel'}"
                            >{line.content}</a>
                        </div>
                    {:else if line.type === 'video-title'}
                        <div class="text-teal-800 font-medium">{line.content}</div>
                    {:else if line.type === 'video-link'}
                        <div class="text-teal-600 text-sm">
                            <a
                                href={line.content.replace('🔗 ', '')}
                                target="_blank"
                                rel="noopener noreferrer"
                                class="hover:underline focus:outline-none focus:ring-2 focus:ring-teal-500 rounded cursor-pointer"
                                tabindex="0"
                                aria-label="Tonton video: {line.attachment?.title || 'Video'}"
                            >{line.content}</a>
                        </div>
                    {:else if line.type === 'attachment'}
                        <div class="text-teal-700">{line.content}</div>
                    {:else}
                        <span>{line.content}</span>
                    {/if}
                {/each}
            </div>

            <!-- Timestamp and scheduled indicator -->
            <div class="flex justify-end items-center gap-1 mt-1">
                {#if isScheduled}
                    <svg class="w-3 h-3 text-amber-500" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M12 2C6.5 2 2 6.5 2 12s4.5 10 10 10 10-4.5 10-10S17.5 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm.5-13H11v6l5.2 3.2.8-1.3-4.5-2.7V7z"/>
                    </svg>
                    <span class="text-[10px] text-amber-600 font-medium">{formattedScheduledTime}</span>
                {:else}
                    <span class="text-[10px] text-slate-500">{formattedTime}</span>
                {/if}
            </div>

            <!-- Tail indicator (sender style) -->
            <div class="absolute -right-2 top-0 w-0 h-0
                border-l-8 border-l-whatsapp-bubble
                border-t-8 border-t-whatsapp-bubble
                border-r-8 border-r-transparent
                border-b-8 border-b-transparent">
            </div>
        </div>
    </div>
</div>

<style>
    /* Additional styling for WhatsApp authenticity */
    .whatsapp-bubble {
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
        box-shadow: 0 1px 0.5px rgba(0, 0, 0, 0.13);
    }
</style>
