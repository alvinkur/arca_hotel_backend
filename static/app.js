// ═══════════════════════════════════════════
// Hotel Arca — Shared Dashboard JS
// ═══════════════════════════════════════════

const App = {
  token: null,
  user: null,

  // ── Auth ──────────────────────────────────

  getToken() {
    return localStorage.getItem('arca_token');
  },

  setToken(t) {
    localStorage.setItem('arca_token', t);
    App.token = t;
  },

  getUser() {
    try { return JSON.parse(localStorage.getItem('arca_user') || 'null'); }
    catch { return null; }
  },

  setUser(u) {
    localStorage.setItem('arca_user', JSON.stringify(u));
    App.user = u;
  },

  logout() {
    localStorage.removeItem('arca_token');
    localStorage.removeItem('arca_user');
    App.token = null;
    App.user = null;
    window.location.href = '/login';
  },

  authHeaders() {
    return {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + App.getToken()
    };
  },

  requireAuth() {
    const t = App.getToken();
    if (!t) {
      window.location.href = '/login';
      return false;
    }
    App.token = t;
    App.user = App.getUser();
    return true;
  },

  // ── API Client ────────────────────────────

  // Safe JSON parse — rejects empty/text responses
  async _json(r) {
    var text = await r.text();
    if (!text) return {};
    try { return JSON.parse(text); }
    catch (e) { throw new Error('Server returned invalid response (not JSON). Pastikan semua service berjalan.'); }
  },

  async _handleError(r) {
    try {
      var data = await App._json(r);
      throw new Error(data.error || 'Request gagal (HTTP ' + r.status + ')');
    } catch (e) {
      if (e.message && e.message.indexOf('Server returned') === 0) throw e;
      throw new Error('Service tidak tersedia (HTTP ' + r.status + '). Pastikan backend berjalan.');
    }
  },

  async list(resource) {
    var r = await fetch('/api/' + resource, { headers: App.authHeaders() });
    if (!r.ok) return App._handleError(r);
    return App._json(r);
  },

  async get(resource, id) {
    var r = await fetch('/api/' + resource + '/' + id, { headers: App.authHeaders() });
    if (!r.ok) return App._handleError(r);
    return App._json(r);
  },

  async create(resource, body) {
    var r = await fetch('/api/' + resource, {
      method: 'POST',
      headers: App.authHeaders(),
      body: JSON.stringify(body)
    });
    if (!r.ok) return App._handleError(r);
    return App._json(r);
  },

  async update(resource, id, body) {
    var r = await fetch('/api/' + resource + '/' + id, {
      method: 'PUT',
      headers: App.authHeaders(),
      body: JSON.stringify(body)
    });
    if (!r.ok) return App._handleError(r);
    return App._json(r);
  },

  async del(resource, id) {
    var r = await fetch('/api/' + resource + '/' + id, {
      method: 'DELETE',
      headers: App.authHeaders()
    });
    if (!r.ok) return App._handleError(r);
    return App._json(r);
  },

  // ── Toast ─────────────────────────────────

  showToast(msg, type) {
    type = type || 'success';
    const t = document.getElementById('toast');
    const icon = document.getElementById('toastIcon');
    t.className = 'toast ' + type;
    document.getElementById('toastMsg').textContent = msg;
    icon.innerHTML = type === 'success'
      ? '<polyline points="20 6 9 17 4 12"/>'
      : '<circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/>';
    t.classList.add('show');
    clearTimeout(App._toastTimer);
    App._toastTimer = setTimeout(function () { t.classList.remove('show'); }, 3000);
  },

  // ── Formatters ────────────────────────────

  initials(name) {
    if (!name) return '?';
    return name.split(' ').map(function (w) { return w[0]; }).join('').substring(0, 2).toUpperCase();
  },

  formatRupiah(n) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(n);
  },

  formatDate(iso) {
    if (!iso) return '—';
    return new Date(iso).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' });
  },

  formatDateTime(iso) {
    if (!iso) return '—';
    return new Date(iso).toLocaleString('id-ID', { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' });
  },

  // ── Modal ─────────────────────────────────

  openModal(title, subtitle, fields, onSubmit, existing) {
    var overlay = document.getElementById('modalOverlay');
    var content = document.getElementById('modalContent');

    existing = existing || {};
    var isEdit = !!existing[Object.keys(existing)[0]];

    var html = '<button class="btn-cancel" onclick="App.closeModal()">' +
      '<svg width="14" height="14" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>' +
      '</button>' +
      '<div class="modal-title">' + title + '</div>' +
      '<div class="modal-sub">' + subtitle + '</div>' +
      '<form id="modalForm" onsubmit="return false">';

    // Build form fields
    var hasGrid = fields.length > 2;
    var gridOpen = false;

    fields.forEach(function (f, i) {
      if (isEdit && f.createOnly) return;

      if (hasGrid && i % 2 === 0 && i < fields.length - 1) {
        html += '<div class="form-grid">';
        gridOpen = true;
      }

      html += '<div class="form-group">';
      html += '<label class="form-label">' + f.label + '</label>';

      var val = existing[f.key] !== undefined ? existing[f.key] : '';

      if (f.type === 'select') {
        html += '<select class="form-select" id="field_' + f.key + '"' + (f.required ? ' required' : '') + '>';
        html += '<option value="">-- Pilih --</option>';
        if (f.options) {
          f.options.forEach(function (opt) {
            var optVal = typeof opt === 'string' ? opt : opt.value;
            var optLabel = typeof opt === 'string' ? opt : (opt.label || opt.value);
            var sel = String(val) === String(optVal) ? ' selected' : '';
            html += '<option value="' + optVal + '"' + sel + '>' + optLabel + '</option>';
          });
        }
        html += '</select>';
      } else if (f.type === 'textarea') {
        html += '<textarea class="form-textarea" id="field_' + f.key + '" placeholder="' + (f.placeholder || '') + '"' + (f.required ? ' required' : '') + '>' + val + '</textarea>';
      } else if (f.type === 'toggle') {
        var checked = val === true || val === 'true' ? ' on' : '';
        html += '<div class="toggle-wrapper">';
        html += '<div class="toggle' + checked + '" id="field_' + f.key + '" onclick="App._toggleSwitch(this)"></div>';
        html += '<span class="toggle-label">' + (f.toggleLabel || '') + '</span>';
        html += '</div>';
      } else {
        var inputType = f.type || 'text';
        html += '<input class="form-input" type="' + inputType + '" id="field_' + f.key + '" placeholder="' + (f.placeholder || '') + '" value="' + App._escAttr(val) + '"' + (f.required ? ' required' : '') + (f.min ? ' min="' + f.min + '"' : '') + (f.max ? ' max="' + f.max + '"' : '') + (f.step ? ' step="' + f.step + '"' : '') + '>';
      }

      html += '</div>';

      if (hasGrid && i % 2 === 1) {
        html += '</div>';
        gridOpen = false;
      }
      if (hasGrid && gridOpen && i === fields.length - 1) {
        html += '<div></div></div>';
        gridOpen = false;
      }
    });

    if (gridOpen) { html += '<div></div></div>'; }

    html += '<button class="btn-submit" type="submit" id="submitBtn" onclick="App._handleSubmit()">' +
      '<svg width="16" height="16" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>' +
      (isEdit ? 'Simpan Perubahan' : 'Simpan') +
      '</button>';

    html += '</form>';

    content.innerHTML = html;
    overlay.classList.add('open');

    // Store callback and edit state
    App._modalOnSubmit = onSubmit;
    App._modalIsEdit = isEdit;
    App._modalExistingId = isEdit ? existing[App._currentEntity.idField] : null;

    // Focus first input
    setTimeout(function () {
      var first = content.querySelector('input:not([type=hidden]), select, textarea');
      if (first) first.focus();
    }, 100);

    // Populate async selects AFTER modal is open
    fields.forEach(function (f) {
      if (f.type === 'select' && f.resource) {
        App._populateSelect('field_' + f.key, f.resource, f.labelField || 'name', f.valueField || App._guessIdField(f.resource), val);
      }
    });
  },

  closeModal() {
    document.getElementById('modalOverlay').classList.remove('open');
    App._modalOnSubmit = null;
    App._modalExistingId = null;
    App._modalIsEdit = false;
  },

  _handleSubmit() {
    var form = document.getElementById('modalForm');
    var fields = App._currentEntity.fields.filter(function (f) { return !(App._modalIsEdit && f.createOnly); });

    var data = {};
    var valid = true;

    fields.forEach(function (f) {
      var el = document.getElementById('field_' + f.key);
      if (!el) return;

      if (f.type === 'toggle') {
        data[f.key] = el.classList.contains('on');
      } else {
        var v = el.value.trim();
        if (f.required && !v) { valid = false; return; }
        if (f.type === 'number') {
          data[f.key] = v ? parseFloat(v) : 0;
        } else if (f.type === 'datetime-local' && v) {
          data[f.key] = new Date(v).toISOString();
        } else {
          data[f.key] = v || '';
        }
      }
    });

    if (!valid) { App.showToast('Mohon lengkapi field yang wajib diisi', 'error'); return; }

    var btn = document.getElementById('submitBtn');
    btn.disabled = true;
    btn.innerHTML = '<div class="spinner"></div> Menyimpan...';

    var promise;
    if (App._modalIsEdit && App._modalExistingId) {
      promise = App._modalOnSubmit(App._modalExistingId, data);
    } else {
      promise = App._modalOnSubmit(data);
    }

    promise
      .then(function () {
        App.closeModal();
        App.loadData(App._currentEntity);
      })
      .catch(function (err) {
        App.showToast(err.message || 'Gagal menyimpan data', 'error');
        btn.disabled = false;
        btn.innerHTML = '<svg width="16" height="16" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg> ' + (App._modalIsEdit ? 'Simpan Perubahan' : 'Simpan');
      });
  },

  // ── Sidebar ───────────────────────────────

  renderSidebar(activePage) {
    var user = App.getUser();
    var container = document.getElementById('sidebarContainer');
    if (!container) return;

    var navItems = [
      { page: 'dashboard', label: 'Dashboard', icon: '<rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/>' },
      { section: 'Manajemen Pengguna' },
      { page: 'customer', label: 'Customers', icon: '<path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/><circle cx="9" cy="7" r="4"/>' },
      { page: 'owner', label: 'Owners', icon: '<path d="M20 21v-2a4 4 0 00-4-4h-4a4 4 0 00-4 4v2"/><circle cx="12" cy="7" r="4"/>' },
      { page: 'staff', label: 'Staffs', icon: '<path d="M16 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/><circle cx="8.5" cy="7" r="4"/><line x1="20" y1="8" x2="20" y2="14"/><line x1="23" y1="11" x2="17" y2="11"/>' },
      { section: 'Manajemen Hotel' },
      { page: 'room', label: 'Rooms', icon: '<path d="M3 9l9-7 9 7v11a2 2 0 01-2 2H5a2 2 0 01-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/>' },
      { page: 'room-type', label: 'Room Types', icon: '<rect x="2" y="7" width="20" height="14" rx="2"/><path d="M16 7V5a2 2 0 00-2-2h-4a2 2 0 00-2 2v2"/>' },
      { page: 'booking', label: 'Bookings', icon: '<path d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>' },
      { page: 'payment', label: 'Payments', icon: '<line x1="12" y1="1" x2="12" y2="23"/><path d="M17 5H9.5a3.5 3.5 0 000 7h5a3.5 3.5 0 010 7H6"/>' },
      { section: 'Lainnya' },
      { page: 'review', label: 'Reviews', icon: '<polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>' },
      { page: 'chat', label: 'Chats', icon: '<path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z"/>' },
      { page: 'revenue-report', label: 'Revenue', icon: '<line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/>' },
    ];

    var html = '<div class="sidebar-logo"><h1>ARCA</h1><p>Hotel Management</p></div><nav>';

    navItems.forEach(function (item) {
      if (item.section) {
        html += '<div class="sidebar-section">' + item.section + '</div>';
      } else {
        var active = activePage === item.page ? ' active' : '';
        html += '<a href="/' + item.page + '" class="nav-item' + active + '">' +
          '<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">' + item.icon + '</svg>' +
          item.label +
          '</a>';
      }
    });

    html += '</nav>';

    // User info at bottom
    if (user) {
      html += '<div class="sidebar-user">' +
        '<div class="avatar">' + App.initials(user.name) + '</div>' +
        '<div class="sidebar-user-info">' +
        '<div class="sidebar-user-name">' + App._escHtml(user.name) + '</div>' +
        '<div class="sidebar-user-role">' + App._escHtml(user.role) + '</div>' +
        '</div>' +
        '<button class="btn-logout" onclick="App.logout()" title="Logout">' +
        '<svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path d="M9 21H5a2 2 0 01-2-2V5a2 2 0 012-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>' +
        '</button>' +
        '</div>';
    }

    container.innerHTML = html;
  },

  // ── Stats ─────────────────────────────────

  renderStats(stats, allData, filteredData) {
    var container = document.getElementById('statsContainer');
    if (!container) return;

    var icons = {
      total: '<path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/><circle cx="9" cy="7" r="4"/>',
      revenue: '<line x1="12" y1="1" x2="12" y2="23"/><path d="M17 5H9.5a3.5 3.5 0 000 7h5a3.5 3.5 0 010 7H6"/>',
      calendar: '<rect x="3" y="4" width="18" height="18" rx="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/>',
      eye: '<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/>',
      star: '<polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>',
      booking: '<path d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>',
    };

    container.innerHTML = stats.map(function (s) {
      return '<div class="stat-card">' +
        '<div class="stat-label">' + s.label + '</div>' +
        '<div class="stat-value" id="stat_' + s.id + '">—</div>' +
        '<div class="stat-icon"><svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">' + (icons[s.icon] || icons.total) + '</svg></div>' +
        '</div>';
    }).join('');

    // Update values
    stats.forEach(function (s) {
      var el = document.getElementById('stat_' + s.id);
      if (el) el.textContent = s.compute(allData, filteredData);
    });
  },

  // ── Table ─────────────────────────────────

  renderTable(data, columns, idField, callbacks) {
    var tbody = document.getElementById('tableBody');
    var thead = document.getElementById('tableHead');

    thead.innerHTML = '<tr>' +
      columns.map(function (c) { return '<th>' + c.label + '</th>'; }).join('') +
      '<th>Aksi</th>' +
      '</tr>';

    if (!data.length) {
      tbody.innerHTML = '<tr><td colspan="' + (columns.length + 1) + '"><div class="empty-state">' +
        '<svg fill="none" viewBox="0 0 24 24" stroke="currentColor"><path d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"/></svg>' +
        '<p>Belum ada data</p></div></td></tr>';
      return;
    }

    tbody.innerHTML = data.map(function (row) {
      var cells = columns.map(function (c) {
        if (c.render) return '<td>' + c.render(row) + '</td>';
        var val = row[c.key];
        if (val === null || val === undefined) val = '—';
        return '<td>' + App._escHtml(String(val)) + '</td>';
      }).join('');

      var actions = '';
      if (callbacks.onEdit) {
        actions += '<button class="btn-edit" onclick="App._currentEntity._onEdit(' + row[idField] + ')">Edit</button>';
      }
      if (callbacks.onDelete) {
        actions += '<button class="btn-delete" onclick="App._currentEntity._onDelete(' + row[idField] + ')">Hapus</button>';
      }

      return '<tr>' + cells + '<td>' + actions + '</td></tr>';
    }).join('');
  },

  // ── Page Init ─────────────────────────────

  _allData: [],
  _currentEntity: null,

  initPage(entity) {
    if (!App.requireAuth()) return;

    App._currentEntity = entity;
    App.user = App.getUser();
    App.renderSidebar(entity.name);

    // _onEdit / _onDelete are set in loadData and _filterTable
    // so that renderTable's onclick handlers reference real callbacks

    // Attach to window for onclick handlers
    window.openCreateModal = function () {
      App.openModal(
        'Tambah ' + entity.displayName,
        'Lengkapi data ' + entity.displayName.toLowerCase() + ' baru',
        entity.fields,
        function (data) { return App.create(entity.name, data); }
      );
    };
    window.openEditModal = function (row) {
      App.openModal(
        'Edit ' + entity.displayName,
        'Perbarui data ' + entity.displayName.toLowerCase(),
        entity.fields,
        function (id, data) { return App.update(entity.name, id, data); },
        row
      );
    };
    window.filterTable = App._filterTable;
    window.handleOverlayClick = function (e) {
      if (e.target === document.getElementById('modalOverlay')) App.closeModal();
    };
    window.deleteEntity = function (id) {
      if (!confirm('Hapus ' + entity.displayName + ' ini?')) return;
      App.del(entity.name, id)
        .then(function () {
          App.showToast(entity.displayName + ' berhasil dihapus!');
          App.loadData(entity);
        })
        .catch(function (err) { App.showToast(err.message, 'error'); });
    };

    App.loadData(entity);
  },

  loadData(entity) {
    App._currentEntity = entity;
    App.list(entity.name)
      .then(function (data) {
        App._allData = data;
        App.renderStats(entity.stats, data, data);

        var callbacks = {
          onEdit: entity.onEdit ? function (id) {
            var row = App._allData.find(function (r) { return r[entity.idField] === id; });
            if (row) window.openEditModal(row);
          } : null,
          onDelete: entity.onDelete !== false ? window.deleteEntity : null
        };

        entity._onEdit = callbacks.onEdit;
        entity._onDelete = callbacks.onDelete;

        App.renderTable(data, entity.columns, entity.idField, callbacks);
      })
      .catch(function (err) {
        App.showToast(err.message || 'Gagal memuat data ' + entity.displayName, 'error');
        document.getElementById('tableBody').innerHTML =
          '<tr><td colspan="' + (entity.columns.length + 1) + '"><div class="empty-state"><p>Gagal memuat data. Periksa koneksi server.</p></div></td></tr>';
      });
  },

  _filterTable() {
    var q = (document.getElementById('searchInput').value || '').toLowerCase();
    var entity = App._currentEntity;

    var filtered = App._allData.filter(function (row) {
      return entity.columns.some(function (c) {
        var val = row[c.key];
        return val !== null && val !== undefined && String(val).toLowerCase().indexOf(q) !== -1;
      });
    });

    App.renderStats(entity.stats, App._allData, filtered);

    var callbacks = {
      onEdit: entity.onEdit ? function (id) {
        var row = filtered.find(function (r) { return r[entity.idField] === id; });
        if (row) window.openEditModal(row);
      } : null,
      onDelete: entity.onDelete !== false ? window.deleteEntity : null
    };

    entity._onEdit = callbacks.onEdit;
    entity._onDelete = callbacks.onDelete;

    App.renderTable(filtered, entity.columns, entity.idField, callbacks);
  },

  // ── Dynamic select populator ──────────────

  _populateSelect(elId, resource, labelField, valueField, selectedVal) {
    var select = document.getElementById(elId);
    if (!select) return;

    App.list(resource)
      .then(function (items) {
        items.forEach(function (item) {
          var opt = document.createElement('option');
          opt.value = item[valueField];
          var label = item[labelField];
          // Append extra info if available
          if (item.price) label += ' (' + App.formatRupiah(item.price) + '/malam)';
          if (item.room_number) label = '#' + item.room_number + ' — ' + label;
          if (item.email) label += ' (' + item.email + ')';
          opt.textContent = label;
          if (String(item[valueField]) === String(selectedVal)) opt.selected = true;
          select.appendChild(opt);
        });
      })
      .catch(function () {});
  },

  // Load select options into an existing select element
  populateSelect(elId, resource, labelField, valueField, placeholder) {
    var select = document.getElementById(elId);
    if (!select) return;
    select.innerHTML = '<option value="">' + (placeholder || '-- Pilih --') + '</option>';
    App.list(resource)
      .then(function (items) {
        items.forEach(function (item) {
          var opt = document.createElement('option');
          opt.value = item[valueField || App._guessIdField(resource)];
          var label = item[labelField || 'name'];
          if (item.price) label += ' (' + App.formatRupiah(item.price) + '/malam)';
          if (item.room_number) label = '#' + item.room_number + ' — ' + label;
          if (item.email) label += ' (' + item.email + ')';
          opt.textContent = label;
          select.appendChild(opt);
        });
      })
      .catch(function (err) { console.error('Gagal load ' + resource + ':', err); });
  },

  _guessIdField(resource) {
    var map = {
      'customers': 'id_customer', 'owners': 'id_owner', 'staffs': 'id_staff',
      'rooms': 'id_room', 'room-types': 'id_room_type', 'bookings': 'id_booking',
      'payments': 'id_payment', 'chats': 'id_chat', 'reviews': 'id_review',
      'revenue_reports': 'id_revenue'
    };
    return map[resource] || 'id';
  },

  // ── Toggle Switch ─────────────────────────

  _toggleSwitch(el) {
    el.classList.toggle('on');
  },

  // ── Escaping helpers ──────────────────────

  _escHtml(s) {
    return String(s).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
  },

  _escAttr(s) {
    return String(s).replace(/&/g, '&amp;').replace(/"/g, '&quot;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
  },

  // ── AI Recommendation ────────────────────

  async askAI(message) {
    var r = await fetch('/ai-recommend', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ message: message })
    });
    var data = await r.json();
    return data.reply;
  }
};
