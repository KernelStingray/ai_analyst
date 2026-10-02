/**
 * SpecAuditor Gov-Edition
 * Frontend Application Engine (Bilingual, File Upload & History Drawer)
 */

document.addEventListener('DOMContentLoaded', () => {
  // Elements
  const docTitleInput = document.getElementById('doc-title');
  const specContentArea = document.getElementById('spec-content');
  const charCountEl = document.getElementById('char-count');
  const wordCountEl = document.getElementById('word-count');
  const charWarningEl = document.getElementById('char-warning');
  
  // Presets & Controls
  const btnLoadPresetUZ = document.getElementById('btn-load-preset-uz');
  const btnLoadPresetRU = document.getElementById('btn-load-preset-ru');
  const btnEmptyLoadPreset = document.getElementById('btn-empty-load-preset');
  const btnClear = document.getElementById('btn-clear');
  const btnAudit = document.getElementById('btn-audit');
  const btnText = document.getElementById('btn-text');

  // File Upload
  const dropzone = document.getElementById('dropzone');
  const fileInput = document.getElementById('file-input');

  // History Drawer
  const btnHistoryToggle = document.getElementById('btn-history-toggle');
  const btnCloseHistory = document.getElementById('btn-close-history');
  const historyDrawer = document.getElementById('history-drawer');
  const historyBackdrop = document.getElementById('history-backdrop');
  const historyList = document.getElementById('history-list');

  // Status & Health
  const statusBadge = document.getElementById('status-badge');
  const statusDot = document.getElementById('status-dot');
  const statusText = document.getElementById('status-text');

  // Panels
  const emptyState = document.getElementById('empty-state');
  const loadingState = document.getElementById('loading-state');
  const loadingTimer = document.getElementById('loading-timer');
  const reportContainer = document.getElementById('report-container');
  const demoModeAlert = document.getElementById('demo-mode-alert');

  // Report Elements
  const reportIdEl = document.getElementById('report-id');
  const reportDateEl = document.getElementById('report-date');
  const reportDocTitleEl = document.getElementById('report-doc-title');
  const scoreCircle = document.getElementById('score-circle');
  const scoreNumberEl = document.getElementById('score-number');
  const verdictBadge = document.getElementById('verdict-badge');
  const verdictSubtext = document.getElementById('verdict-subtext');
  const metricStruct = document.getElementById('metric-struct');
  const metricScenarios = document.getElementById('metric-scenarios');
  const metricMeasurability = document.getElementById('metric-measurability');
  const executiveSummaryText = document.getElementById('executive-summary-text');
  const missingSectionsWrapper = document.getElementById('missing-sections-wrapper');
  const missingSectionsList = document.getElementById('missing-sections-list');
  const findingsCountEl = document.getElementById('findings-count');
  const findingsList = document.getElementById('findings-list');
  const btnPrintPdf = document.getElementById('btn-print-pdf');

  // Print elements
  const printDocTitle = document.getElementById('print-doc-title');
  const printDocDate = document.getElementById('print-doc-date');
  const printDocId = document.getElementById('print-doc-id');

  let currentFindings = [];
  let currentFilter = 'ALL';
  let timerInterval = null;
  const MAX_CHARS = 500000;

  // 1. Initial Health Check
  checkHealth();
  setInterval(checkHealth, 15000);

  async function checkHealth() {
    try {
      const res = await fetch('/api/health');
      const data = await res.json();
      if (data.success && data.data) {
        const info = data.data;
        if (info.llm_available) {
          statusBadge.className = 'flex items-center space-x-2 px-3 py-1.5 rounded-lg border text-xs font-mono transition-all duration-300 bg-emerald-500/10 border-emerald-500/30 text-emerald-400';
          statusDot.innerHTML = `
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-2.5 w-2.5 bg-emerald-500"></span>
          `;
          statusText.textContent = `Lokal LLM: ${info.llm_model} (Faol)`;
        } else {
          statusBadge.className = 'flex items-center space-x-2 px-3 py-1.5 rounded-lg border text-xs font-mono transition-all duration-300 bg-amber-500/10 border-amber-500/30 text-amber-300';
          statusDot.innerHTML = `
            <span class="relative inline-flex rounded-full h-2.5 w-2.5 bg-amber-400"></span>
          `;
          statusText.textContent = 'Rejim: Offline Demo Koʻrgazma';
        }
      }
    } catch (e) {
      console.warn('Health check failed:', e);
      statusBadge.className = 'flex items-center space-x-2 px-3 py-1.5 rounded-lg border text-xs font-mono transition-all duration-300 bg-amber-500/10 border-amber-500/30 text-amber-300';
      statusText.textContent = 'Rejim: Offline Demo Koʻrgazma';
    }
  }

  // 2. Character & Word Counter
  function updateCounters() {
    const text = specContentArea.value || '';
    const charLen = text.length;
    const words = text.trim() ? text.trim().split(/\s+/).length : 0;

    charCountEl.textContent = `${charLen.toLocaleString()} / ${MAX_CHARS.toLocaleString()} belgi`;
    wordCountEl.textContent = `${words.toLocaleString()} soʻz`;

    if (charLen > MAX_CHARS) {
      charCountEl.classList.add('text-rose-400', 'font-bold');
      charWarningEl.classList.remove('hidden');
      btnAudit.disabled = true;
    } else {
      charCountEl.classList.remove('text-rose-400', 'font-bold');
      charWarningEl.classList.add('hidden');
      btnAudit.disabled = false;
    }
  }

  specContentArea.addEventListener('input', updateCounters);

  // 3. Preset Loading (Bilingual)
  async function loadPreset(lang = 'uz') {
    try {
      if (btnLoadPresetUZ) btnLoadPresetUZ.disabled = true;
      if (btnLoadPresetRU) btnLoadPresetRU.disabled = true;

      const res = await fetch(`/api/preset?lang=${lang}`);
      const data = await res.json();
      if (data.success && data.data && data.data.preset) {
        specContentArea.value = data.data.preset;
        if (lang === 'ru') {
          docTitleInput.value = 'АИС Мониторинг госуслуг ТЗ (Проект)';
        } else {
          docTitleInput.value = 'DXM-Stat AT Texnik topshirigʻi (Namunaviy loyiha)';
        }
        updateCounters();
        specContentArea.focus();
      }
    } catch (err) {
      alert('Namunaviy matnni yuklashda xatolik: ' + err.message);
    } finally {
      if (btnLoadPresetUZ) btnLoadPresetUZ.disabled = false;
      if (btnLoadPresetRU) btnLoadPresetRU.disabled = false;
    }
  }

  if (btnLoadPresetUZ) {
    btnLoadPresetUZ.addEventListener('click', () => loadPreset('uz'));
  }
  if (btnLoadPresetRU) {
    btnLoadPresetRU.addEventListener('click', () => loadPreset('ru'));
  }
  if (btnEmptyLoadPreset) {
    btnEmptyLoadPreset.addEventListener('click', () => loadPreset('uz'));
  }

  // 4. File Drag & Drop + Upload
  if (dropzone && fileInput) {
    dropzone.addEventListener('click', () => fileInput.click());

    dropzone.addEventListener('dragover', (e) => {
      e.preventDefault();
      dropzone.classList.add('border-indigo-500', 'bg-indigo-950/30');
    });

    ['dragleave', 'dragend'].forEach(type => {
      dropzone.addEventListener(type, () => {
        dropzone.classList.remove('border-indigo-500', 'bg-indigo-950/30');
      });
    });

    dropzone.addEventListener('drop', (e) => {
      e.preventDefault();
      dropzone.classList.remove('border-indigo-500', 'bg-indigo-950/30');
      if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
        handleFileUpload(e.dataTransfer.files[0]);
      }
    });

    fileInput.addEventListener('change', (e) => {
      if (e.target.files && e.target.files.length > 0) {
        handleFileUpload(e.target.files[0]);
      }
    });
  }

  async function handleFileUpload(file) {
    const formData = new FormData();
    formData.append('file', file);

    const originalText = dropzone.innerHTML;
    dropzone.innerHTML = `
      <div class="flex items-center space-x-2 text-indigo-400 text-xs py-1">
        <svg class="w-4 h-4 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m0 14v1m8-8h-1M5 12H4m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707" />
        </svg>
        <span>Fayl ochilmoqda (${escapeHtml(file.name)})...</span>
      </div>
    `;

    try {
      const res = await fetch('/api/upload', {
        method: 'POST',
        body: formData,
      });
      const data = await res.json();
      if (!res.ok || !data.success) {
        throw new Error(data.error || 'Faylni oʻqishda xatolik');
      }

      specContentArea.value = data.data.content || '';
      docTitleInput.value = data.data.document_name || file.name;
      updateCounters();
    } catch (err) {
      alert('Fayl xatosi: ' + err.message);
    } finally {
      dropzone.innerHTML = originalText;
      fileInput.value = '';
    }
  }

  // 5. History Drawer
  if (btnHistoryToggle && historyDrawer && historyBackdrop) {
    btnHistoryToggle.addEventListener('click', openHistory);
    if (btnCloseHistory) btnCloseHistory.addEventListener('click', closeHistory);
    historyBackdrop.addEventListener('click', closeHistory);
  }

  function openHistory() {
    historyDrawer.classList.remove('translate-x-full');
    historyBackdrop.classList.remove('hidden');
    loadHistory();
  }

  function closeHistory() {
    historyDrawer.classList.add('translate-x-full');
    historyBackdrop.classList.add('hidden');
  }

  async function loadHistory() {
    try {
      historyList.innerHTML = `
        <div class="text-center py-6 text-xs text-slate-400">
          <svg class="w-5 h-5 mx-auto mb-2 animate-spin text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Tarix yuklanmoqda...
        </div>
      `;
      const res = await fetch('/api/history?limit=15');
      if (!res.ok) {
        throw new Error(`Server javobi: ${res.status} (aloqa uzildi yoki server ishga tushmoqda)`);
      }
      const ctype = res.headers.get('content-type') || '';
      if (!ctype.includes('application/json')) {
        throw new Error('Serverdan kutilmagan javob qaytdi. Sahifani yangilang (F5).');
      }
      const data = await res.json();
      if (!data.success || !data.data || data.data.length === 0) {
        historyList.innerHTML = `<div class="text-center py-8 text-xs text-slate-500">Hozircha hech qanday ekspertiza oʻtkazilmagan.</div>`;
        return;
      }

      historyList.innerHTML = '';
      data.data.forEach(item => {
        const card = document.createElement('div');
        const isApproved = item.status === 'TASDIQLANDI';
        const isConditional = item.status === 'SHARTLI TASDIQLANDI';
        const badgeColor = isApproved ? 'text-emerald-400 bg-emerald-500/10 border-emerald-500/30' : (isConditional ? 'text-amber-400 bg-amber-500/10 border-amber-500/30' : 'text-rose-400 bg-rose-500/10 border-rose-500/30');

        card.className = 'p-3 rounded-xl bg-slate-800/80 hover:bg-slate-800 border border-slate-700/80 hover:border-indigo-500/60 transition-all cursor-pointer space-y-1.5';
        card.innerHTML = `
          <div class="flex items-center justify-between">
            <span class="text-[10px] font-mono font-bold px-2 py-0.5 rounded border ${badgeColor}">${item.status || 'RAD ETILDI'}</span>
            <span class="text-xs font-mono font-bold text-white">${Math.round((item.compliance_score || 0)*10)/10}%</span>
          </div>
          <h4 class="text-xs font-semibold text-slate-200 line-clamp-1">${escapeHtml(item.document_name || 'Nomsiz TT')}</h4>
          <div class="flex items-center justify-between text-[10px] text-slate-400 font-mono">
            <span>#${item.audit_id || ''}</span>
            <span>${new Date(item.analyzed_at).toLocaleTimeString('uz-UZ', {hour: '2-digit', minute:'2-digit'})}</span>
          </div>
        `;
        card.addEventListener('click', () => {
          renderReport(item);
          closeHistory();
        });
        historyList.appendChild(card);
      });
    } catch (e) {
      historyList.innerHTML = `<div class="text-center py-4 text-xs text-rose-400 leading-relaxed">${e.message}</div>`;
    }
  }

  // 6. Clear Field
  btnClear.addEventListener('click', () => {
    if (confirm('Maydonni tozalashni xohlaysizmi?')) {
      specContentArea.value = '';
      docTitleInput.value = '';
      updateCounters();
      emptyState.classList.remove('hidden');
      reportContainer.classList.add('hidden');
      loadingState.classList.add('hidden');
    }
  });

  // 7. Submit Audit
  btnAudit.addEventListener('click', async () => {
    const content = specContentArea.value.trim();
    if (!content) {
      alert('Iltimos, avval Texnik topshiriq matnini kiriting yoki "Namunaviy TT" tugmasini bosing.');
      return;
    }

    if (content.length < 50) {
      alert('Matn hajmi juda qisqa (kamida 50 ta belgi boʻlishi lozim).');
      return;
    }

    if (content.length > MAX_CHARS) {
      alert(`Matn hajmi ${MAX_CHARS} belgidan oshmasligi kerak.`);
      return;
    }

    // Switch to loading state
    startLoading();

    try {
      const response = await fetch('/api/audit', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          document_name: docTitleInput.value.trim() || 'Nomsiz TT Loyihasi',
          content: content,
        }),
      });

      const resData = await response.json();
      if (!response.ok || !resData.success) {
        throw new Error(resData.error || 'Ekspertiza jarayonida xatolik yuz berdi');
      }

      renderReport(resData.data);
    } catch (err) {
      alert('Xatolik: ' + err.message);
      stopLoading(false);
    }
  });

  let pipelineInterval = null;

  function updateStep(stepNum, status) {
    const el = document.getElementById(`step-${stepNum}`);
    if (!el) return;
    const icon = el.querySelector('.step-icon');

    if (status === 'active') {
      el.className = 'flex items-center space-x-3 text-xs text-indigo-300 transition-all duration-300 font-semibold';
      if (icon) {
        icon.className = 'step-icon w-5 h-5 rounded-full bg-indigo-500/20 border border-indigo-500 flex items-center justify-center flex-shrink-0';
        icon.innerHTML = '<span class="w-2 h-2 rounded-full bg-indigo-400 animate-ping"></span>';
      }
    } else if (status === 'done') {
      el.className = 'flex items-center space-x-3 text-xs text-emerald-400 transition-all duration-300';
      if (icon) {
        icon.className = 'step-icon w-5 h-5 rounded-full bg-emerald-500/20 border border-emerald-500/50 flex items-center justify-center flex-shrink-0 text-emerald-400';
        icon.innerHTML = '<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7"/></svg>';
      }
    } else {
      el.className = 'flex items-center space-x-3 text-xs text-slate-400 transition-all duration-300';
      if (icon) {
        icon.className = 'step-icon w-5 h-5 rounded-full bg-slate-700/80 flex items-center justify-center flex-shrink-0';
        icon.innerHTML = '<span class="w-2 h-2 rounded-full bg-slate-500"></span>';
      }
    }
  }

  function startLoading() {
    emptyState.classList.add('hidden');
    reportContainer.classList.add('hidden');
    loadingState.classList.remove('hidden');

    btnAudit.disabled = true;
    btnText.textContent = 'TAHLIL QILINMOQDA...';

    // Reset steps to pending
    for (let i = 1; i <= 5; i++) updateStep(i, 'pending');
    updateStep(1, 'active');

    const progressBar = document.getElementById('loading-progress-bar');
    if (progressBar) progressBar.style.width = '15%';

    let currentStep = 1;
    let seconds = 0;
    loadingTimer.textContent = 'Hujjat standart boʻyicha boʻlaklarga ajratilmoqda (0 soniya)...';

    clearInterval(timerInterval);
    timerInterval = setInterval(() => {
      seconds++;
      loadingTimer.textContent = `Avtomatlashtirilgan ekspertiza jarayoni (${seconds} soniya)...`;
    }, 1000);

    clearInterval(pipelineInterval);
    pipelineInterval = setInterval(() => {
      if (currentStep < 5) {
        updateStep(currentStep, 'done');
        currentStep++;
        updateStep(currentStep, 'active');
        if (progressBar) {
          const widths = [15, 35, 60, 80, 95];
          progressBar.style.width = `${widths[currentStep - 1]}%`;
        }
      }
    }, 1100);
  }

  function stopLoading(showReport = true) {
    clearInterval(timerInterval);
    clearInterval(pipelineInterval);

    for (let i = 1; i <= 5; i++) updateStep(i, 'done');
    const progressBar = document.getElementById('loading-progress-bar');
    if (progressBar) progressBar.style.width = '100%';

    loadingState.classList.add('hidden');
    btnAudit.disabled = false;
    btnText.textContent = 'EKSPERTIZADAN OʻTKAZISH';

    if (showReport) {
      reportContainer.classList.remove('hidden');
      reportContainer.scrollIntoView({ behavior: 'smooth' });
    } else {
      emptyState.classList.remove('hidden');
    }
  }

  // 8. Render Audit Report
  function renderReport(report) {
    stopLoading(true);

    // Metadata
    reportIdEl.textContent = '#' + (report.audit_id || 'aud-001');
    const dateObj = new Date(report.analyzed_at || Date.now());
    const dateFormatted = dateObj.toLocaleDateString('uz-UZ') + ', ' + dateObj.toLocaleTimeString('uz-UZ', { hour: '2-digit', minute: '2-digit' });
    reportDateEl.textContent = dateFormatted;
    reportDocTitleEl.textContent = report.document_name || 'Texnik topshiriq loyihasi';

    // Print headers
    printDocTitle.textContent = report.document_name || '-';
    printDocDate.textContent = dateFormatted;
    printDocId.textContent = report.audit_id || '-';

    // Offline demo badge
    if (report.is_demo_mode) {
      demoModeAlert.classList.remove('hidden');
    } else {
      demoModeAlert.classList.add('hidden');
    }

    // Score & Gauge (Circle circumference = 2 * PI * 40 ~= 251.2)
    const score = Math.round((report.compliance_score || 0) * 10) / 10;
    scoreNumberEl.textContent = `${score}%`;

    const circumference = 251.2;
    const offset = circumference - (score / 100) * circumference;
    scoreCircle.style.strokeDashoffset = offset;

    // Status / Verdict
    const status = report.status || (score >= 85 ? 'TASDIQLANDI' : (score >= 70 ? 'SHARTLI TASDIQLANDI' : 'RAD ETILDI'));
    verdictBadge.textContent = status;

    if (status === 'TASDIQLANDI') {
      verdictBadge.className = 'px-3 py-1 rounded-md font-bold text-xs font-mono uppercase bg-emerald-500/20 text-emerald-300 border border-emerald-500/30';
      scoreCircle.className.baseVal = 'text-emerald-500 transition-all duration-1000';
      verdictSubtext.textContent = 'Hujjat Oʻz DSt 1987:2018 talablariga toʻliq javob beradi va davlat ekspertizasidan oʻtishi mumkin.';
    } else if (status === 'SHARTLI TASDIQLANDI') {
      verdictBadge.className = 'px-3 py-1 rounded-md font-bold text-xs font-mono uppercase bg-amber-500/20 text-amber-300 border border-amber-500/30';
      scoreCircle.className.baseVal = 'text-amber-500 transition-all duration-1000';
      verdictSubtext.textContent = 'Qayd etilgan kamchiliklar bartaraf etilgandan soʻng takroriy ekspertizaga taqdim etilsin.';
    } else {
      verdictBadge.className = 'px-3 py-1 rounded-md font-bold text-xs font-mono uppercase bg-rose-500/20 text-rose-300 border border-rose-500/30';
      scoreCircle.className.baseVal = 'text-rose-500 transition-all duration-1000';
      verdictSubtext.textContent = 'Standart talablari qoʻpol ravishda buzilgan. Hujjat tender yoki ishlab chiqishga kiritilishi taqiqlanadi.';
    }

    // Metrics Cards
    if (report.metrics) {
      metricStruct.textContent = `${report.metrics.structure_completeness || 0}%`;
      metricScenarios.textContent = `${report.metrics.scenarios_compliance || 0}%`;
      metricMeasurability.textContent = `${report.metrics.measurability_score || 0}%`;
    }

    // Executive Summary
    executiveSummaryText.textContent = report.executive_summary || 'Tahlil yakunlandi.';

    // Section Breakdown Grid (Auto-slicing visualization)
    const sectionsWrapper = document.getElementById('sections-map-wrapper');
    const sectionsGrid = document.getElementById('sections-grid');
    const detectedCountEl = document.getElementById('detected-sections-count');

    if (sectionsWrapper && sectionsGrid && report.section_map && report.section_map.length > 0) {
      sectionsWrapper.classList.remove('hidden');
      sectionsGrid.innerHTML = '';
      const detectedTotal = report.section_map.filter(s => s.detected).length;
      if (detectedCountEl) {
        detectedCountEl.textContent = `${detectedTotal} / ${report.section_map.length} boʻlim aniqlandi`;
      }

      report.section_map.forEach(sec => {
        const card = document.createElement('div');
        const isDetected = sec.detected;
        const issues = sec.issues_count || 0;

        let borderClass = 'border-slate-800 bg-slate-800/40';
        let badgeHtml = '';

        if (!isDetected) {
          borderClass = 'border-rose-500/30 bg-rose-500/5';
          badgeHtml = `<span class="px-2 py-0.5 rounded text-[10px] font-mono font-bold uppercase bg-rose-500/20 text-rose-300 border border-rose-500/30">TUSHIRIB QOLDIRILGAN</span>`;
        } else if (issues > 0) {
          borderClass = 'border-amber-500/30 bg-amber-500/5';
          badgeHtml = `<span class="px-2 py-0.5 rounded text-[10px] font-mono font-bold uppercase bg-amber-500/20 text-amber-300 border border-amber-500/30">${issues} TA NUQSON</span>`;
        } else {
          borderClass = 'border-emerald-500/30 bg-emerald-500/5';
          badgeHtml = `<span class="px-2 py-0.5 rounded text-[10px] font-mono font-bold uppercase bg-emerald-500/20 text-emerald-300 border border-emerald-500/30">MUVOFИQ</span>`;
        }

        card.className = `p-3.5 rounded-xl border ${borderClass} space-y-2 transition-all`;
        card.innerHTML = `
          <div class="flex items-start justify-between gap-2">
            <div>
              <span class="text-xs font-bold text-white block">${escapeHtml(sec.title)}</span>
              <span class="text-[10px] font-mono text-slate-400">${escapeHtml(sec.standard_ref || '')}</span>
            </div>
            ${badgeHtml}
          </div>
          ${sec.snippet ? `
            <div class="text-[11px] text-slate-300 font-mono line-clamp-2 bg-slate-900/60 p-2 rounded border border-slate-800/80">
              «${escapeHtml(sec.snippet)}»
            </div>
          ` : `
            <div class="text-[11px] text-rose-400 italic">
              Boʻlim matni hujjatdan topilmadi (Oʻz DSt 1987:2018 talabi buzilgan)
            </div>
          `}
        `;
        sectionsGrid.appendChild(card);
      });
    } else if (sectionsWrapper) {
      sectionsWrapper.classList.add('hidden');
    }

    // Missing Sections
    missingSectionsList.innerHTML = '';
    if (report.missing_sections && report.missing_sections.length > 0) {
      missingSectionsWrapper.classList.remove('hidden');
      report.missing_sections.forEach(sec => {
        const item = document.createElement('div');
        item.className = 'p-3 rounded-xl bg-rose-500/10 border border-rose-500/20 text-xs text-rose-200 flex items-start space-x-2.5';
        item.innerHTML = `
          <svg class="w-4 h-4 text-rose-400 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
          <div class="font-medium">${escapeHtml(sec)}</div>
        `;
        missingSectionsList.appendChild(item);
      });
    } else {
      missingSectionsWrapper.classList.add('hidden');
    }

    // Findings Feed
    currentFindings = report.findings || [];
    findingsCountEl.textContent = currentFindings.length;
    renderFindingsFeed();
  }

  // 9. Filter & Render Findings Feed
  function renderFindingsFeed() {
    findingsList.innerHTML = '';

    const filtered = currentFindings.filter(f => {
      if (currentFilter === 'ALL') return true;
      return (f.severity || '').toUpperCase() === currentFilter;
    });

    if (filtered.length === 0) {
      findingsList.innerHTML = `
        <div class="text-center py-6 text-xs text-slate-500 border border-slate-800 rounded-xl">
          Tanlangan filtr boʻyicha kamchiliklar mavjud emas.
        </div>
      `;
      return;
    }

    filtered.forEach((finding, idx) => {
      const card = document.createElement('div');
      card.className = 'bg-gov-card border border-gov-border rounded-xl p-4 space-y-3 hover:border-slate-600 transition-all';

      const isCritical = finding.severity === 'CRITICAL';
      const isWarning = finding.severity === 'WARNING';
      const badgeClass = isCritical
        ? 'bg-rose-500/20 text-rose-300 border-rose-500/30'
        : (isWarning ? 'bg-amber-500/20 text-amber-300 border-amber-500/30' : 'bg-blue-500/20 text-blue-300 border-blue-500/30');

      card.innerHTML = `
        <div class="flex items-start justify-between gap-2">
          <div class="flex items-center space-x-2">
            <span class="px-2 py-0.5 rounded text-[10px] font-mono font-bold uppercase border ${badgeClass}">
              ${finding.severity || 'NOTE'}
            </span>
            <span class="text-xs font-semibold text-white">${escapeHtml(finding.section || 'Boʻlim')}</span>
          </div>
          <span class="text-[11px] font-mono text-slate-500">#${finding.id || 'FIND-' + (idx + 1)}</span>
        </div>

        <!-- Raw Fragment with error -->
        ${finding.raw_fragment ? `
          <div class="p-2.5 rounded-lg bg-slate-900/90 border-l-4 border-rose-500 text-xs font-mono text-slate-300 leading-relaxed">
            <span class="text-slate-500 text-[10px] block mb-0.5 uppercase tracking-wider font-sans font-bold">TTdagi nuqsonli matn:</span>
            «${escapeHtml(finding.raw_fragment)}»
          </div>
        ` : ''}

        <!-- Violation Clause -->
        <div class="text-xs text-slate-300 flex items-start space-x-2">
          <svg class="w-3.5 h-3.5 text-rose-400 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01" />
          </svg>
          <div>
            <strong class="text-slate-400">Tafovut:</strong> ${escapeHtml(finding.violation || '')}
          </div>
        </div>

        <!-- Recommendation / Ready text -->
        ${finding.recommendation ? `
          <div class="p-3 rounded-lg bg-emerald-950/30 border border-emerald-500/30 text-xs space-y-1.5">
            <div class="flex items-center justify-between text-emerald-400 font-semibold text-[11px]">
              <span class="flex items-center">
                <svg class="w-3.5 h-3.5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                Tavsiya etilgan toʻgʻri tahrir:
              </span>
              <button class="btn-copy-rec text-[11px] text-emerald-300 hover:text-white px-2 py-0.5 rounded bg-emerald-500/20 hover:bg-emerald-500/40 transition-all" data-rec="${escapeHtmlAttr(finding.recommendation)}">
                Nusxa olish
              </button>
            </div>
            <p class="text-slate-200 font-sans leading-relaxed">${escapeHtml(finding.recommendation)}</p>
          </div>
        ` : ''}
      `;

      findingsList.appendChild(card);
    });

    // Wire up Copy Recommendation buttons
    document.querySelectorAll('.btn-copy-rec').forEach(btn => {
      btn.addEventListener('click', (e) => {
        const textToCopy = e.target.getAttribute('data-rec');
        if (textToCopy) {
          navigator.clipboard.writeText(textToCopy).then(() => {
            const orig = e.target.textContent;
            e.target.textContent = 'Nusxalandi!';
            setTimeout(() => { e.target.textContent = orig; }, 1800);
          });
        }
      });
    });
  }

  // 10. Filter buttons handler
  document.querySelectorAll('.btn-filter').forEach(btn => {
    btn.addEventListener('click', () => {
      document.querySelectorAll('.btn-filter').forEach(b => {
        b.className = 'btn-filter px-2.5 py-1 rounded bg-slate-800/60 text-slate-400 hover:text-slate-200 font-medium';
      });
      btn.className = 'btn-filter px-2.5 py-1 rounded bg-indigo-600 text-white font-medium shadow';
      currentFilter = btn.getAttribute('data-filter');
      renderFindingsFeed();
    });
  });

  // 11. PDF Print Trigger
  btnPrintPdf.addEventListener('click', () => {
    window.print();
  });

  // Utility helpers
  function escapeHtml(str) {
    if (!str) return '';
    return str
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#039;');
  }

  function escapeHtmlAttr(str) {
    if (!str) return '';
    return str.replace(/"/g, '&quot;');
  }
});
