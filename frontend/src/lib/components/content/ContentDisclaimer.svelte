<script>
  import { t } from 'svelte-i18n';

  let { class: className = "" } = $props();

  // Fallback disclaimer text if i18n key is missing
  const fallbackDisclaimer = "Konten ini untuk tujuan edukasi kesehatan. Konsultasikan dengan tenaga kesehatan untuk kondisi spesifik Anda.";

  let disclaimerText = $derived($t('content.disclaimer.text') || fallbackDisclaimer);
</script>

<div
  class="mt-4 pt-3 border-t border-slate-100 {className}"
  role="region"
  aria-label="Health disclaimer"
>
  <p class="text-xs text-slate-500 leading-relaxed">
    {disclaimerText}
  </p>
</div>
