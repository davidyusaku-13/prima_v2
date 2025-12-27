<script>
  import { t, locale } from 'svelte-i18n';

  export let activities = [];

  function formatTime(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    const now = new Date();
    const diff = now - date;

    if (diff < 60000) return $t('activity.justNow');
    if (diff < 3600000) return $t('activity.minutesAgo', { n: Math.floor(diff / 60000) });
    if (diff < 86400000) return $t('activity.hoursAgo', { n: Math.floor(diff / 3600000) });
    if (diff < 604800000) return $t('activity.daysAgo', { n: Math.floor(diff / 86400000) });

    return date.toLocaleDateString($locale || 'en-US', { day: 'numeric', month: 'short' });
  }

  function getActivityIcon(type) {
    const icons = {
      article_created: '<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" /></svg>',
      article_updated: '<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>',
      article_deleted: '<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>',
      video_created: '<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" /></svg>',
      video_deleted: '<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>',
      user_login: '<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1" /></svg>'
    };
    return icons[type] || icons.article_created;
  }

  function getActivityColor(type) {
    const colors = {
      article_created: 'bg-green-100 text-green-600',
      article_updated: 'bg-blue-100 text-blue-600',
      article_deleted: 'bg-red-100 text-red-600',
      video_created: 'bg-purple-100 text-purple-600',
      video_deleted: 'bg-red-100 text-red-600',
      user_login: 'bg-teal-100 text-teal-600'
    };
    return colors[type] || 'bg-slate-100 text-slate-600';
  }

  function getActivityMessage(activity) {
    const messages = {
      article_created: $t('activity.articleCreated', { title: activity.title }),
      article_updated: $t('activity.articleUpdated', { title: activity.title }),
      article_deleted: $t('activity.articleDeleted', { title: activity.title }),
      video_created: $t('activity.videoCreated', { title: activity.title }),
      video_deleted: $t('activity.videoDeleted'),
      user_login: $t('activity.userLogin')
    };
    return messages[activity.type] || activity.type;
  }
</script>

<div class="bg-white rounded-2xl border border-slate-200">
  <div class="px-5 py-4 border-b border-slate-100">
    <h3 class="font-semibold text-slate-900">{$t('cms.recentActivity')}</h3>
  </div>

  {#if activities.length === 0}
    <div class="px-5 py-8 text-center">
      <svg class="w-12 h-12 mx-auto text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <p class="mt-3 text-sm text-slate-500">{$t('cms.noActivity')}</p>
    </div>
  {:else}
    <div class="divide-y divide-slate-100">
      {#each activities as activity}
        <div class="px-5 py-3 flex items-center gap-3 hover:bg-slate-50 transition-colors">
          <div class="w-8 h-8 rounded-full flex items-center justify-center {getActivityColor(activity.type)}">
            {@html getActivityIcon(activity.type)}
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm text-slate-900 truncate">{getActivityMessage(activity)}</p>
            <p class="text-xs text-slate-400">{activity.user || $t('common.system')}</p>
          </div>
          <span class="text-xs text-slate-400 whitespace-nowrap">
            {formatTime(activity.createdAt)}
          </span>
        </div>
      {/each}
    </div>
  {/if}
</div>
