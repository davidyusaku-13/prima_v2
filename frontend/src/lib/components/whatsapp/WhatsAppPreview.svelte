<script>
    import { t } from 'svelte-i18n';

    let {
        message = '',
        patientName = '',
        reminderTitle = '',
        timestamp = new Date(),
        showDisclaimer = true
    } = $props();

    // Format timestamp to HH:mm
    let formattedTime = $derived.by(() => {
        const date = timestamp instanceof Date ? timestamp : new Date(timestamp);
        return date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' });
    });

    // Greeting text (i18n)
    let greetingText = $derived($t('reminder.preview.greeting', { values: { name: patientName } }));

    // Disclaimer text (i18n)
    let disclaimerText = $derived($t('reminder.preview.disclaimer'));

    // Build complete message preview - returns array of line objects for rendering
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

        // Disclaimer
        if (showDisclaimer) {
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
                    {:else}
                        <span>{line.content}</span>
                    {/if}
                {/each}
            </div>

            <!-- Timestamp -->
            <div class="flex justify-end items-center gap-1 mt-1">
                <span class="text-[10px] text-slate-500">{formattedTime}</span>
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
