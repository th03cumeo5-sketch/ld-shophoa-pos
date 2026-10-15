/**
 * ShopHoa POS - Frontend Application Engine
 * Integrates with Wails Go Backend & SQL Server 2019
 * Designed for high-speed cashier workflow with keyboard shortcuts.
 */

(function () {
  'use strict';

  // =========================================================================
  // STATE MANAGEMENT
  // =========================================================================
  const state = {
    isWails: Boolean(window.go && window.go.main && window.go.main.App),
    session: {
      makho: '01',
      quay: '01',
      ca: 'Ca Sáng',
      thungan: 'ntd',
      loc: 'Cửa hàng Q2 (Kho: K01)',
      time: ''
    },
    currentSale: {
      header: {
        khoa: '',
        lct: 'BLE',
        sct: '00001',
        ngay: new Date().toISOString(),
        loc: '01',
        makho: '01',
        quay: '01',
        ca: '1',
        thungan: 'ntd',
        madt: '',
        mavip: '',
        soluong: 0,
        sotien: 0,
        chietkhau_mh: 0,
        tl_ck: 0,
        chietkhau: 0,
        thanhtien_chua_vat: 0,
        tien_thue: 0,
        tien_thue_cl: 0,
        thanhtoan: 0,
        pttt: 'TM',
        ttoan1: 0,
        ttoan2: 0,
        chuathoi: 0,
        printed: false,
        print_no: 0,
        dgiai: ''
      },
      items: []
    },
    customer: {
      madt: '',
      tendt: 'Khách lẻ',
      dienthoai: '',
      mavip: ''
    },
    selectedItemIndex: -1,
    selectedCategory: '',
    catalogItems: [],
    categories: [],
    draftSalesCount: 0,
    isReturnMode: false,
    searchResults: []
  };

  // Backend API Proxy (uses Wails runtime if available, else mock engine)
  const api = {
    async getSessionInfo() {
      if (state.isWails) {
        return window.go.main.App.GetSessionInfo();
      }
      return state.session;
    },
    async prepareNewSale() {
      if (state.isWails) {
        return window.go.main.App.PrepareNewSale();
      }
      // Mock bill generator
      const mockSct = String(Math.floor(Math.random() * 90000) + 10000);
      return {
        khoa: 'mock-' + Date.now(),
        lct: 'BLE',
        sct: mockSct,
        ngay: new Date().toISOString(),
        loc: '01',
        makho: '01',
        quay: '01',
        ca: '1',
        thungan: 'ntd',
        printed: false
      };
    },
    async searchProducts(keyword, category) {
      if (state.isWails) {
        return window.go.main.App.SearchProducts(keyword, category || '');
      }
      return mockProducts.filter(p => {
        const matchKw = !keyword || p.tenvt.toLowerCase().includes(keyword.toLowerCase()) || p.mavt.includes(keyword);
        const matchCat = !category || p.manhom === category;
        return matchKw && matchCat;
      });
    },
    async getCategories() {
      if (state.isWails) {
        return window.go.main.App.GetCategories();
      }
      return [
        { manhom: '2402', tennhom: 'Hoa tươi', count: 120 },
        { manhom: '2401', tennhom: 'Ly + Bình hoa', count: 45 },
        { manhom: '1901', tennhom: 'Lưu niệm & Quà', count: 68 },
        { manhom: '3101', tennhom: 'Thiệp & Nơ', count: 35 },
        { manhom: '5005', tennhom: 'Giấy gói quà', count: 50 },
        { manhom: '2403', tennhom: 'Giỏ trái cây', count: 25 }
      ];
    },
    async getCatalogProducts(limit) {
      if (state.isWails) {
        return window.go.main.App.GetCatalogProducts(limit || 24);
      }
      return mockProducts.slice(0, limit || 24);
    },
    async calculateCart(sale) {
      if (state.isWails) {
        return window.go.main.App.CalculateCart(sale);
      }
      return localCalculateCart(sale);
    },
    async saveDraftSale(sale) {
      if (state.isWails) {
        return window.go.main.App.SaveDraftSale(sale);
      }
      mockDrafts.push(JSON.parse(JSON.stringify(sale)));
      return true;
    },
    async completePayment(sale, pttt, ttoan1, ttoan2) {
      if (state.isWails) {
        return window.go.main.App.CompletePayment(sale, pttt, ttoan1, ttoan2);
      }
      return true;
    },
    async listDraftSales() {
      if (state.isWails) {
        return window.go.main.App.ListDraftSales();
      }
      return mockDrafts.map(d => d.header);
    },
    async getDraftSale(khoa) {
      if (state.isWails) {
        return window.go.main.App.GetDraftSale(khoa);
      }
      return mockDrafts.find(d => d.header.khoa === khoa) || null;
    },
    async deleteDraftSale(khoa) {
      if (state.isWails) {
        return window.go.main.App.DeleteDraftSale(khoa);
      }
      const idx = mockDrafts.findIndex(d => d.header.khoa === khoa);
      if (idx >= 0) mockDrafts.splice(idx, 1);
      return true;
    },
    async searchCustomers(keyword) {
      if (state.isWails) {
        return window.go.main.App.SearchCustomers(keyword);
      }
      return mockCustomers.filter(c => !keyword || c.tendt.toLowerCase().includes(keyword.toLowerCase()) || c.dienthoai.includes(keyword));
    }
  };

  // Mock Data for offline preview & testing
  const mockProducts = [
    { mavt: 'HOA001', tenvt: 'Hoa Hồng Đỏ Ecuador cao cấp', dvt: 'Cành', giaban: 35000, loaithue: 'T10', vat_ra: 10, manhom: '2402' },
    { mavt: 'HOA002', tenvt: 'Hoa Hồng Ohara Trắng thơm', dvt: 'Cành', giaban: 40000, loaithue: 'T10', vat_ra: 10, manhom: '2402' },
    { mavt: 'HOA003', tenvt: 'Cúc Tana Hà Lan nhập khẩu', dvt: 'Bó', giaban: 120000, loaithue: 'T05', vat_ra: 5, manhom: '2402' },
    { mavt: 'HOA004', tenvt: 'Baby Trắng phun màu Pastel', dvt: 'Bó', giaban: 150000, loaithue: 'T05', vat_ra: 5, manhom: '2402' },
    { mavt: 'HOA005', tenvt: 'Lan Hồ Điệp Tím 5 cành', dvt: 'Chậu', giaban: 1250000, loaithue: 'T10', vat_ra: 10, manhom: '2402' },
    { mavt: 'HOA006', tenvt: 'Lẵng hoa Khai Trương Phát Tài', dvt: 'Lẵng', giaban: 850000, loaithue: 'T10', vat_ra: 10, manhom: '2402' },
    { mavt: 'BINH01', tenvt: 'Bình hoa Thủy tinh Pha lê', dvt: 'Cái', giaban: 220000, loaithue: 'T10', vat_ra: 10, manhom: '2401' },
    { mavt: 'BINH02', tenvt: 'Bình gốm sứ Bát Tràng cắm hoa', dvt: 'Cái', giaban: 180000, loaithue: 'T10', vat_ra: 10, manhom: '2401' },
    { mavt: 'GIFT01', tenvt: 'Gấu bông tốt nghiệp mini', dvt: 'Con', giaban: 95000, loaithue: 'T10', vat_ra: 10, manhom: '1901' },
    { mavt: 'THIEP1', tenvt: 'Thiệp hoa khô Handmade', dvt: 'Cái', giaban: 25000, loaithue: 'T08', vat_ra: 8, manhom: '3101' },
    { mavt: 'GOI001', tenvt: 'Giấy gói hoa chống thấm Hàn Quốc', dvt: 'Xấp', giaban: 65000, loaithue: 'T08', vat_ra: 8, manhom: '5005' },
    { mavt: 'GOI002', tenvt: 'Ruy băng lụa thắt nơ cao cấp', dvt: 'Cuộn', giaban: 45000, loaithue: 'T08', vat_ra: 8, manhom: '5005' }
  ];

  const mockCustomers = [
    { madt: 'KH001', tendt: 'Nguyễn Thị Mai', dienthoai: '0903123456', mavip: 'VIP01' },
    { madt: 'KH002', tendt: 'Trần Văn Hùng', dienthoai: '0918765432', mavip: '' },
    { madt: 'KH003', tendt: 'Công ty Cổ phần Alpha', dienthoai: '0283888999', mavip: 'VIP02' }
  ];

  let mockDrafts = [];

  // =========================================================================
  // DOM ELEMENTS
  // =========================================================================
  const dom = {
    barcodeInput: document.getElementById('barcode-input'),
    searchResultsDropdown: document.getElementById('search-results-dropdown'),
    orderSct: document.getElementById('order-sct'),
    orderTime: document.getElementById('order-time'),
    customerName: document.getElementById('customer-name'),
    btnCustomer: document.getElementById('btn-customer'),
    cartItemsBody: document.getElementById('cart-items-body'),
    cartEmpty: document.getElementById('cart-empty'),
    totalQty: document.getElementById('total-qty'),
    subtotalVal: document.getElementById('subtotal-val'),
    discountVal: document.getElementById('discount-val'),
    taxVal: document.getElementById('tax-val'),
    grandTotalVal: document.getElementById('grand-total-val'),
    catalogTabs: document.getElementById('catalog-tabs'),
    catalogGrid: document.getElementById('catalog-grid'),
    catalogFilterInput: document.getElementById('catalog-filter-input'),
    draftCountBadge: document.getElementById('draft-count'),
    btnChangeQty: document.getElementById('btn-change-qty'),
    btnDiscount: document.getElementById('btn-discount'),
    btnDiscountTag: document.getElementById('btn-discount-tag'),
    btnHoldOrder: document.getElementById('btn-hold-order'),
    btnPendingDrafts: document.getElementById('btn-pending-drafts'),
    btnReturnMode: document.getElementById('btn-return-mode'),
    btnCancelOrder: document.getElementById('btn-cancel-order'),
    btnCheckout: document.getElementById('btn-checkout'),
    btnSearchOrder: document.getElementById('btn-search-order'),
    btnReprint: document.getElementById('btn-reprint'),
    btnLock: document.getElementById('btn-lock'),
    sessionShift: document.getElementById('session-shift'),
    sessionCashier: document.getElementById('session-cashier'),
    sessionStore: document.getElementById('session-store'),
    sessionClock: document.getElementById('session-clock'),
    toastContainer: document.getElementById('toast-container'),

    // Modals
    modalPayment: document.getElementById('modal-payment'),
    modalPaySct: document.getElementById('modal-pay-sct'),
    modalPayTotal: document.getElementById('modal-pay-total'),
    inputCashGiven: document.getElementById('input-cash-given'),
    modalChangeDue: document.getElementById('modal-change-due'),
    btnSubmitPay: document.getElementById('btn-submit-pay'),
    btnCancelPay: document.getElementById('btn-cancel-pay'),
    sectionCash: document.getElementById('section-cash'),
    sectionQr: document.getElementById('section-qr'),
    qrMemo: document.getElementById('qr-memo'),

    modalQty: document.getElementById('modal-qty'),
    qtyItemTitle: document.getElementById('qty-item-title'),
    inputNewQty: document.getElementById('input-new-qty'),
    btnSubmitQty: document.getElementById('btn-submit-qty'),

    modalDiscount: document.getElementById('modal-discount'),
    discountInputLabel: document.getElementById('discount-input-label'),
    inputDiscountVal: document.getElementById('input-discount-val'),
    btnSubmitDiscount: document.getElementById('btn-submit-discount'),

    modalDrafts: document.getElementById('modal-drafts'),
    draftSalesBody: document.getElementById('draft-sales-body'),

    modalCustomer: document.getElementById('modal-customer'),
    inputSearchCustomer: document.getElementById('input-search-customer'),
    customerListBox: document.getElementById('customer-list-box'),
    btnChooseWalkin: document.getElementById('btn-choose-walkin'),

    modalLock: document.getElementById('modal-lock'),
    inputUnlockPass: document.getElementById('input-unlock-pass'),
    btnUnlock: document.getElementById('btn-unlock')
  };

  // =========================================================================
  // UTILITY FUNCTIONS
  // =========================================================================
  function formatMoney(amount) {
    if (isNaN(amount) || amount === null) amount = 0;
    return new Intl.NumberFormat('vi-VN').format(Math.round(amount)) + ' đ';
  }

  function parseMoney(str) {
    if (!str) return 0;
    const clean = String(str).replace(/[^\d.-]/g, '');
    return parseFloat(clean) || 0;
  }

  function showToast(message, type = 'info') {
    const toast = document.createElement('div');
    toast.className = 'toast';
    if (type === 'success') toast.style.backgroundColor = 'var(--btn-green)';
    if (type === 'error') toast.style.backgroundColor = 'var(--btn-danger)';
    if (type === 'warning') toast.style.backgroundColor = 'var(--btn-amber)';
    toast.textContent = message;
    dom.toastContainer.appendChild(toast);
    setTimeout(() => {
      toast.style.opacity = '0';
      toast.style.transition = 'opacity 0.3s';
      setTimeout(() => toast.remove(), 300);
    }, 2800);
  }

  function updateClock() {
    const now = new Date();
    dom.sessionClock.textContent = now.toTimeString().split(' ')[0];
  }

  // =========================================================================
  // CORE CALCULATION ENGINE (Offline Mirror of Go Core Service)
  // =========================================================================
  function localCalculateCart(sale) {
    let totalQty = 0;
    let subtotal = 0;
    let totalCKLines = 0;
    let totalTax = 0;

    sale.items.forEach((item, idx) => {
      item.stt = idx + 1;
      const qty = parseFloat(item.soluong) || 0;
      const price = parseFloat(item.dongia) || 0;
      const lineTotal = qty * price;
      item.sotien = lineTotal;

      // Line discount rate: tổng hợp các tỷ lệ chiết khấu
      const totalLineRate = (parseFloat(item.tl_ck2) || 0) + (parseFloat(item.tl_ck3) || 0) +
                            (parseFloat(item.tl_ck4) || 0) + (parseFloat(item.tl_ck5) || 0) + (parseFloat(item.tl_ck) || 0);
      const rate = Math.min(totalLineRate, 50.0); // trần 50%
      const lineCK = (lineTotal * rate) / 100.0;
      item.chietkhau = lineCK;
      item.dongia_sauck = qty > 0 ? (lineTotal - lineCK) / qty : price;

      // Tax
      const taxRate = parseFloat(item.vat_ra) || 10;
      const taxable = lineTotal - lineCK;
      const lineTax = (taxable * taxRate) / 100.0;
      item.thue_suat = taxRate;
      item.tien_thue = lineTax;
      item.thanhtien_chua_vat = taxable;
      item.thanhtoan = taxable + lineTax;

      totalQty += qty;
      subtotal += lineTotal;
      totalCKLines += lineCK;
      totalTax += lineTax;
    });

    // Header discount
    let headerCK = 0;
    if (sale.header.tl_ck > 0) {
      headerCK = (subtotal * sale.header.tl_ck) / 100.0;
    } else if (sale.header.chietkhau > 0) {
      headerCK = sale.header.chietkhau;
    }
    sale.header.chietkhau = headerCK;
    sale.header.chietkhau_mh = totalCKLines;

    sale.header.soluong = totalQty;
    sale.header.sotien = subtotal;
    sale.header.thanhtien_chua_vat = subtotal - totalCKLines - headerCK;
    sale.header.tien_thue = totalTax;
    sale.header.thanhtoan = Math.max(0, (subtotal - totalCKLines - headerCK) + totalTax);

    return sale;
  }

  // =========================================================================
  // UI RENDERING
  // =========================================================================
  function renderCart() {
    const { items, header } = state.currentSale;
    dom.cartItemsBody.innerHTML = '';

    if (items.length === 0) {
      dom.cartEmpty.style.display = 'block';
      state.selectedItemIndex = -1;
    } else {
      dom.cartEmpty.style.display = 'none';
      if (state.selectedItemIndex < 0 || state.selectedItemIndex >= items.length) {
        state.selectedItemIndex = items.length - 1;
      }

      items.forEach((item, index) => {
        const tr = document.createElement('tr');
        tr.className = 'cart-row' + (index === state.selectedItemIndex ? ' selected' : '');
        tr.dataset.index = index;

        tr.innerHTML = `
          <td style="text-align: center; font-weight: 700; color: var(--text-muted);">${item.stt}</td>
          <td>
            <div class="cart-item-name">${item.tenvt}</div>
            <div class="cart-item-code">${item.mavt} ${item.ghichu ? ' • ' + item.ghichu : ''}</div>
          </td>
          <td style="text-align: center; color: var(--text-secondary);">${item.dvt}</td>
          <td style="text-align: center;">
            <span class="qty-badge" data-action="change-qty" data-index="${index}">${item.soluong}</span>
          </td>
          <td style="text-align: right;" class="cart-price">${formatMoney(item.dongia)}</td>
          <td style="text-align: center; font-weight: 600; color: ${item.chietkhau > 0 ? 'var(--btn-amber)' : 'var(--text-muted)'};">
            ${item.tl_ck ? item.tl_ck + '%' : '-'}
          </td>
          <td style="text-align: right;" class="cart-amount">${formatMoney(item.thanhtoan)}</td>
          <td style="text-align: center;">
            <button class="btn btn-rose btn-sm" data-action="delete-item" data-index="${index}" title="Xóa dòng này">✕</button>
          </td>
        `;

        tr.addEventListener('click', (e) => {
          if (e.target.dataset.action === 'delete-item') {
            removeCartItem(index);
          } else if (e.target.dataset.action === 'change-qty') {
            state.selectedItemIndex = index;
            openQtyModal();
          } else {
            state.selectedItemIndex = index;
            renderCartHighlight();
          }
        });

        dom.cartItemsBody.appendChild(tr);
      });
    }

    // Update Summary box
    dom.orderSct.textContent = header.sct || '---';
    dom.totalQty.textContent = `${header.soluong || 0} món`;
    dom.subtotalVal.textContent = formatMoney(header.sotien);
    dom.discountVal.textContent = `- ${formatMoney((header.chietkhau || 0) + (header.chietkhau_mh || 0))}`;
    dom.taxVal.textContent = formatMoney(header.tien_thue);
    dom.grandTotalVal.innerHTML = `${formatMoney(header.thanhtoan)}`;
    dom.customerName.textContent = state.customer.tendt || 'Khách lẻ';
  }

  function renderCartHighlight() {
    const rows = dom.cartItemsBody.querySelectorAll('.cart-row');
    rows.forEach((row, idx) => {
      if (idx === state.selectedItemIndex) {
        row.classList.add('selected');
        row.scrollIntoView({ block: 'nearest', behavior: 'smooth' });
      } else {
        row.classList.remove('selected');
      }
    });
  }

  function renderCategories() {
    dom.catalogTabs.innerHTML = '';

    const allBtn = document.createElement('button');
    allBtn.className = 'catalog-tab' + (state.selectedCategory === '' ? ' active' : '');
    allBtn.textContent = 'Tất cả';
    allBtn.addEventListener('click', () => {
      state.selectedCategory = '';
      renderCategories();
      loadCatalogProducts();
    });
    dom.catalogTabs.appendChild(allBtn);

    state.categories.forEach(cat => {
      const btn = document.createElement('button');
      btn.className = 'catalog-tab' + (state.selectedCategory === cat.manhom ? ' active' : '');
      btn.textContent = `${cat.tennhom} (${cat.count || 0})`;
      btn.addEventListener('click', () => {
        state.selectedCategory = cat.manhom;
        renderCategories();
        loadCatalogProducts();
      });
      dom.catalogTabs.appendChild(btn);
    });
  }

  function renderCatalogGrid() {
    dom.catalogGrid.innerHTML = '';
    const kw = dom.catalogFilterInput.value.trim().toLowerCase();

    const filtered = state.catalogItems.filter(p => {
      if (!kw) return true;
      return p.tenvt.toLowerCase().includes(kw) || p.mavt.toLowerCase().includes(kw);
    });

    if (filtered.length === 0) {
      dom.catalogGrid.innerHTML = `
        <div style="grid-column: 1 / -1; text-align: center; padding: 40px; color: var(--text-muted); font-size: 13px;">
          Không có mặt hàng nào phù hợp với bộ lọc.
        </div>
      `;
      return;
    }

    filtered.forEach(p => {
      const card = document.createElement('div');
      card.className = 'catalog-card';
      card.innerHTML = `
        <div>
          <div class="catalog-card-name" title="${p.tenvt}">${p.tenvt}</div>
          <div class="catalog-card-meta">Mã: ${p.mavt} • ĐVT: ${p.dvt}</div>
        </div>
        <div class="catalog-card-price">${formatMoney(p.giaban)}</div>
      `;

      card.addEventListener('click', () => {
        addProductToCart(p);
      });

      dom.catalogGrid.appendChild(card);
    });
  }

  // =========================================================================
  // CART OPERATIONS
  // =========================================================================
  async function addProductToCart(product, qty = 1) {
    if (state.isReturnMode) {
      qty = -Math.abs(qty);
    }

    const existingIdx = state.currentSale.items.findIndex(i => i.mavt === product.mavt);

    if (existingIdx >= 0) {
      // Tăng số lượng dòng đã có
      state.currentSale.items[existingIdx].soluong += qty;
      state.selectedItemIndex = existingIdx;
      showToast(`Đã tăng [${product.tenvt}] lên SL: ${state.currentSale.items[existingIdx].soluong}`, 'info');
    } else {
      // Thêm dòng mới
      const newItem = {
        khoact: 'ct-' + Date.now() + '-' + Math.random().toString(36).substr(2, 4),
        stt: state.currentSale.items.length + 1,
        mavt: product.mavt,
        tenvt: product.tenvt,
        dvt: product.dvt || 'Cái',
        soluong: qty,
        dongia: product.giaban || 0,
        dongia_sauck: product.giaban || 0,
        sotien: (product.giaban || 0) * qty,
        tl_ck: 0,
        tl_ck2: state.customer.mavip ? 5 : 0, // 5% nếu là khách VIP
        tl_ck3: 0,
        tl_ck4: 0,
        tl_ck5: 0,
        chietkhau: 0,
        loaithue: product.loaithue || 'T10',
        thue_suat: product.vat_ra || 10,
        tien_thue: 0,
        tien_thue_cl: 0,
        thanhtien_chua_vat: 0,
        thanhtoan: 0,
        ghichu: state.isReturnMode ? 'Khách trả hàng' : ''
      };

      state.currentSale.items.push(newItem);
      state.selectedItemIndex = state.currentSale.items.length - 1;
      showToast(`Đã thêm [${product.tenvt}] vào đơn`, 'success');
    }

    await recalculateAndRefresh();
    dom.barcodeInput.value = '';
    dom.barcodeInput.focus();
  }

  async function removeCartItem(index) {
    if (index >= 0 && index < state.currentSale.items.length) {
      const removed = state.currentSale.items.splice(index, 1)[0];
      if (state.selectedItemIndex >= state.currentSale.items.length) {
        state.selectedItemIndex = state.currentSale.items.length - 1;
      }
      showToast(`Đã xóa [${removed.tenvt}] khỏi đơn`, 'info');
      await recalculateAndRefresh();
    }
  }

  async function recalculateAndRefresh() {
    const updated = await api.calculateCart(state.currentSale);
    state.currentSale = updated;
    renderCart();
  }

  async function startNewSale() {
    const fresh = await api.prepareNewSale();
    state.currentSale = {
      header: {
        ...state.currentSale.header,
        khoa: fresh.khoa,
        sct: fresh.sct,
        ngay: fresh.ngay || new Date().toISOString(),
        soluong: 0,
        sotien: 0,
        chietkhau: 0,
        thanhtien_chua_vat: 0,
        tien_thue: 0,
        thanhtoan: 0,
        printed: false
      },
      items: []
    };
    state.selectedItemIndex = -1;
    state.isReturnMode = false;
    dom.btnReturnMode.style.backgroundColor = 'var(--btn-rose)';
    dom.btnReturnMode.textContent = '[F8] Trả hàng';

    await recalculateAndRefresh();
    await updateDraftCount();
    dom.barcodeInput.focus();
  }

  async function updateDraftCount() {
    try {
      const drafts = await api.listDraftSales();
      state.draftSalesCount = drafts ? drafts.length : 0;
      dom.draftCountBadge.textContent = state.draftSalesCount;
    } catch (e) {
      console.warn('Could not update draft count', e);
    }
  }

  // =========================================================================
  // SEARCH / BARCODE ENGINE
  // =========================================================================
  let searchDebounceTimer = null;

  async function handleSearchInput(e) {
    const val = dom.barcodeInput.value.trim();

    if (e.key === 'Enter') {
      if (!val) return;
      // Tìm chính xác mã barcode hoặc mã vật tư
      const results = await api.searchProducts(val, '');
      if (results && results.length > 0) {
        addProductToCart(results[0]);
        closeSearchDropdown();
      } else {
        showToast(`Không tìm thấy sản phẩm có mã [${val}]`, 'warning');
      }
      return;
    }

    if (e.key === 'Escape') {
      closeSearchDropdown();
      return;
    }

    clearTimeout(searchDebounceTimer);
    if (!val || val.length < 2) {
      closeSearchDropdown();
      return;
    }

    searchDebounceTimer = setTimeout(async () => {
      const results = await api.searchProducts(val, '');
      state.searchResults = results || [];
      renderSearchDropdown();
    }, 200);
  }

  function renderSearchDropdown() {
    if (!state.searchResults || state.searchResults.length === 0) {
      closeSearchDropdown();
      return;
    }

    dom.searchResultsDropdown.innerHTML = '';
    state.searchResults.forEach((p, idx) => {
      const item = document.createElement('div');
      item.className = 'search-item';
      item.innerHTML = `
        <div>
          <div class="search-item-title">${p.tenvt}</div>
          <div class="search-item-meta">Mã: ${p.mavt} • ĐVT: ${p.dvt}</div>
        </div>
        <div class="search-item-price">${formatMoney(p.giaban)}</div>
      `;

      item.addEventListener('click', () => {
        addProductToCart(p);
        closeSearchDropdown();
      });

      dom.searchResultsDropdown.appendChild(item);
    });

    dom.searchResultsDropdown.classList.add('open');
  }

  function closeSearchDropdown() {
    dom.searchResultsDropdown.classList.remove('open');
    dom.searchResultsDropdown.innerHTML = '';
  }

  // =========================================================================
  // MODALS LOGIC
  // =========================================================================
  function openModal(modal) {
    modal.classList.add('open');
  }

  function closeModal(modal) {
    modal.classList.remove('open');
    dom.barcodeInput.focus();
  }

  // 1. PAYMENT MODAL
  function openPaymentModal() {
    if (state.currentSale.items.length === 0) {
      showToast('Giỏ hàng chưa có sản phẩm để thanh toán!', 'warning');
      return;
    }

    const total = state.currentSale.header.thanhtoan;
    dom.modalPaySct.textContent = state.currentSale.header.sct;
    dom.modalPayTotal.textContent = formatMoney(total);
    dom.inputCashGiven.value = Math.round(total).toLocaleString('vi-VN');
    dom.qrMemo.textContent = 'BLE-' + state.currentSale.header.sct;
    updatePaymentChange();

    // Reset payment method to Cash (TM)
    selectPaymentMethod('TM');

    openModal(dom.modalPayment);
    dom.inputCashGiven.select();
  }

  function selectPaymentMethod(method) {
    state.currentSale.header.pttt = method;
    dom.modalPayment.querySelectorAll('.payment-method-btn').forEach(btn => {
      if (btn.dataset.method === method) {
        btn.classList.add('active');
      } else {
        btn.classList.remove('active');
      }
    });

    if (method === 'TM') {
      dom.sectionCash.style.display = 'flex';
      dom.sectionQr.style.display = 'none';
      dom.inputCashGiven.focus();
    } else if (method === 'CK') {
      dom.sectionCash.style.display = 'none';
      dom.sectionQr.style.display = 'block';
    } else {
      dom.sectionCash.style.display = 'none';
      dom.sectionQr.style.display = 'none';
    }
  }

  function updatePaymentChange() {
    const required = state.currentSale.header.thanhtoan;
    const given = parseMoney(dom.inputCashGiven.value);
    const diff = given - required;

    if (diff >= 0) {
      dom.modalChangeDue.className = 'change-due-val positive';
      dom.modalChangeDue.textContent = formatMoney(diff);
      dom.btnSubmitPay.disabled = false;
      dom.btnSubmitPay.style.opacity = '1';
    } else {
      dom.modalChangeDue.className = 'change-due-val negative';
      dom.modalChangeDue.textContent = `Thiếu ${formatMoney(Math.abs(diff))}`;
      dom.btnSubmitPay.disabled = false;
    }
  }

  async function executePayment() {
    const required = state.currentSale.header.thanhtoan;
    const given = parseMoney(dom.inputCashGiven.value);
    const method = state.currentSale.header.pttt || 'TM';

    if (method === 'TM' && given < required) {
      showToast('Số tiền khách đưa chưa đủ để thanh toán!', 'error');
      dom.inputCashGiven.focus();
      return;
    }

    try {
      showToast('Đang tiến hành lưu hóa đơn và thanh toán POS_PAID...', 'info');
      dom.btnSubmitPay.textContent = 'Đang xử lý...';
      dom.btnSubmitPay.disabled = true;

      const ttoan1 = method === 'TM' ? given : required;
      await api.completePayment(state.currentSale, method, ttoan1, 0);

      closeModal(dom.modalPayment);
      showToast(`✓ Đã thanh toán thành công hóa đơn [${state.currentSale.header.sct}] - Tổng: ${formatMoney(required)}`, 'success');

      // Khởi tạo ngay hóa đơn mới cho lượt khách kế tiếp
      await startNewSale();
    } catch (err) {
      console.error(err);
      showToast(`Lỗi thanh toán: ${err.message || err}`, 'error');
    } finally {
      dom.btnSubmitPay.textContent = '[Enter] XÁC NHẬN THANH TOÁN & IN';
      dom.btnSubmitPay.disabled = false;
    }
  }

  // 2. QUANTITY MODAL
  function openQtyModal() {
    const idx = state.selectedItemIndex;
    if (idx < 0 || idx >= state.currentSale.items.length) {
      showToast('Vui lòng chọn 1 dòng sản phẩm trên giỏ hàng trước!', 'warning');
      return;
    }

    const item = state.currentSale.items[idx];
    dom.qtyItemTitle.textContent = `${item.mavt} - ${item.tenvt} (Hiện tại: ${item.soluong} ${item.dvt})`;
    dom.inputNewQty.value = item.soluong;
    openModal(dom.modalQty);
    dom.inputNewQty.select();
  }

  async function submitQtyModal() {
    const newQty = parseFloat(dom.inputNewQty.value) || 1;
    const idx = state.selectedItemIndex;
    if (idx >= 0 && idx < state.currentSale.items.length) {
      state.currentSale.items[idx].soluong = newQty;
      closeModal(dom.modalQty);
      await recalculateAndRefresh();
      showToast(`Đã cập nhật số lượng dòng [${state.currentSale.items[idx].tenvt}] thành ${newQty}`, 'success');
    }
  }

  // 3. DISCOUNT MODAL
  function openDiscountModal() {
    openModal(dom.modalDiscount);
    dom.inputDiscountVal.value = state.currentSale.header.tl_ck || '';
    dom.inputDiscountVal.focus();
    dom.inputDiscountVal.select();
  }

  async function submitDiscountModal() {
    const type = dom.modalDiscount.querySelector('input[name="discount-type"]:checked').value;
    const val = parseFloat(dom.inputDiscountVal.value) || 0;

    if (type === 'percent') {
      state.currentSale.header.tl_ck = val;
      state.currentSale.header.chietkhau = 0;
    } else {
      state.currentSale.header.tl_ck = 0;
      state.currentSale.header.chietkhau = val;
    }

    closeModal(dom.modalDiscount);
    await recalculateAndRefresh();
    showToast(`Đã áp dụng giảm giá đơn hàng: ${val}${type === 'percent' ? '%' : ' đ'}`, 'success');
  }

  // 4. HOLD ORDER (LƯU TẠM BANLE_TAM)
  async function holdCurrentOrder() {
    if (state.currentSale.items.length === 0) {
      showToast('Không thể giữ đơn vì giỏ hàng đang trống!', 'warning');
      return;
    }

    try {
      showToast(`Đang lưu đơn tạm [${state.currentSale.header.sct}] vào BANLE_TAM...`, 'info');
      await api.saveDraftSale(state.currentSale);
      showToast(`✓ Đã lưu tạm đơn hàng [${state.currentSale.header.sct}] thành công!`, 'success');
      await startNewSale();
    } catch (err) {
      console.error(err);
      showToast(`Lỗi khi lưu đơn tạm: ${err.message || err}`, 'error');
    }
  }

  // 5. LIST DRAFT SALES MODAL
  async function openDraftsModal() {
    try {
      const drafts = await api.listDraftSales();
      dom.draftSalesBody.innerHTML = '';

      if (!drafts || drafts.length === 0) {
        dom.draftSalesBody.innerHTML = `
          <tr><td colspan="6" style="text-align: center; padding: 30px; color: var(--text-muted);">
            Hiện không có đơn tạm nào đang treo tại quầy này.
          </td></tr>
        `;
      } else {
        drafts.forEach(d => {
          const tr = document.createElement('tr');
          const timeStr = d.ngay ? new Date(d.ngay).toLocaleTimeString('vi-VN') : '';
          tr.innerHTML = `
            <td><b style="color: var(--color-primary);">${d.sct}</b></td>
            <td>${timeStr}</td>
            <td>${d.madt || 'Khách lẻ'}</td>
            <td style="text-align: center;"><b>${d.soluong}</b></td>
            <td style="text-align: right; font-weight: 800; color: var(--btn-green);">${formatMoney(d.thanhtoan)}</td>
            <td style="text-align: center;">
              <button class="btn btn-blue btn-sm" data-action="resume-draft" data-khoa="${d.khoa}">Mở lại</button>
              <button class="btn btn-rose btn-sm" data-action="del-draft" data-khoa="${d.khoa}">Xóa</button>
            </td>
          `;

          tr.querySelector('[data-action="resume-draft"]').addEventListener('click', async () => {
            await resumeDraftSale(d.khoa);
          });
          tr.querySelector('[data-action="del-draft"]').addEventListener('click', async () => {
            if (confirm(`Bạn có chắc muốn xóa vĩnh viễn đơn tạm [${d.sct}]?`)) {
              await api.deleteDraftSale(d.khoa);
              showToast(`Đã xóa đơn tạm [${d.sct}]`, 'info');
              await updateDraftCount();
              openDraftsModal();
            }
          });

          dom.draftSalesBody.appendChild(tr);
        });
      }

      openModal(dom.modalDrafts);
    } catch (err) {
      console.error(err);
      showToast(`Lỗi nạp danh sách đơn tạm: ${err.message || err}`, 'error');
    }
  }

  async function resumeDraftSale(khoa) {
    try {
      const fullDraft = await api.getDraftSale(khoa);
      if (fullDraft) {
        state.currentSale = fullDraft;
        state.customer.madt = fullDraft.header.madt || '';
        state.customer.tendt = fullDraft.header.madt ? 'Khách ' + fullDraft.header.madt : 'Khách lẻ';
        closeModal(dom.modalDrafts);
        renderCart();
        await updateDraftCount();
        showToast(`✓ Đã phục hồi đơn tạm [${fullDraft.header.sct}] lên màn hình bán hàng`, 'success');
      }
    } catch (err) {
      console.error(err);
      showToast(`Lỗi khi mở lại đơn tạm: ${err.message || err}`, 'error');
    }
  }

  // 6. CUSTOMER SEARCH MODAL
  function openCustomerModal() {
    openModal(dom.modalCustomer);
    dom.inputSearchCustomer.value = '';
    renderCustomerList([]);
    dom.inputSearchCustomer.focus();
  }

  async function searchCustomers() {
    const kw = dom.inputSearchCustomer.value.trim();
    const list = await api.searchCustomers(kw);
    renderCustomerList(list || []);
  }

  function renderCustomerList(list) {
    dom.customerListBox.innerHTML = '';
    if (!list || list.length === 0) {
      dom.customerListBox.innerHTML = `
        <div style="padding: 20px; text-align: center; color: var(--text-muted); font-size: 12px;">
          Gõ số điện thoại hoặc tên khách hàng để tìm kiếm...
        </div>
      `;
      return;
    }

    list.forEach(c => {
      const item = document.createElement('div');
      item.className = 'search-item';
      item.innerHTML = `
        <div>
          <div class="search-item-title">${c.tendt} ${c.mavip ? '<span class="badge" style="background:var(--btn-amber); color:#fff; font-size:10px; padding:2px 4px; border-radius:3px;">VIP</span>' : ''}</div>
          <div class="search-item-meta">Mã: ${c.madt} • SĐT: ${c.dienthoai || 'Không có'}</div>
        </div>
        <button class="btn btn-blue btn-sm">Chọn</button>
      `;

      item.addEventListener('click', () => {
        setCustomer(c);
        closeModal(dom.modalCustomer);
      });

      dom.customerListBox.appendChild(item);
    });
  }

  function setCustomer(c) {
    state.customer = c;
    state.currentSale.header.madt = c.madt;
    state.currentSale.header.mavip = c.mavip || '';
    dom.customerName.textContent = c.tendt;
    showToast(`Đã chọn khách hàng: ${c.tendt}`, 'info');
    recalculateAndRefresh();
  }

  // =========================================================================
  // KEYBOARD SHORTCUTS ENGINE
  // =========================================================================
  function handleGlobalKeyDown(e) {
    const anyModalOpen = document.querySelector('.modal-overlay.open');

    // Hotkey: Escape
    if (e.key === 'Escape') {
      if (anyModalOpen) {
        closeModal(anyModalOpen);
        e.preventDefault();
        return;
      }
      if (dom.barcodeInput.value) {
        dom.barcodeInput.value = '';
        closeSearchDropdown();
        e.preventDefault();
        return;
      }
      if (state.currentSale.items.length > 0) {
        if (confirm('Bạn có chắc chắn muốn hủy đơn hàng đang soạn thảo này không?')) {
          startNewSale();
        }
        e.preventDefault();
      }
      return;
    }

    // Hotkey: Enter in Modals
    if (e.key === 'Enter') {
      if (dom.modalPayment.classList.contains('open')) {
        executePayment();
        e.preventDefault();
        return;
      }
      if (dom.modalQty.classList.contains('open')) {
        submitQtyModal();
        e.preventDefault();
        return;
      }
      if (dom.modalDiscount.classList.contains('open')) {
        submitDiscountModal();
        e.preventDefault();
        return;
      }
    }

    // In modals, skip POS global function keys
    if (anyModalOpen) {
      if (e.key === '1' && dom.modalPayment.classList.contains('open') && document.activeElement !== dom.inputCashGiven) {
        selectPaymentMethod('TM');
        e.preventDefault();
      }
      if (e.key === '2' && dom.modalPayment.classList.contains('open') && document.activeElement !== dom.inputCashGiven) {
        selectPaymentMethod('CK');
        e.preventDefault();
      }
      if (e.key === '3' && dom.modalPayment.classList.contains('open') && document.activeElement !== dom.inputCashGiven) {
        selectPaymentMethod('THE');
        e.preventDefault();
      }
      return;
    }

    // POS Global Function Keys (Strictly aligned with legacy & MISA conventions)
    switch (e.key) {
      case 'F2':
        e.preventDefault();
        dom.barcodeInput.focus();
        dom.barcodeInput.select();
        break;

      case 'F3':
        e.preventDefault();
        openQtyModal();
        break;

      case 'F4':
        e.preventDefault();
        openPaymentModal();
        break;

      case 'F5':
        e.preventDefault();
        openDiscountModal();
        break;

      case 'F6':
        e.preventDefault();
        openCustomerModal();
        break;

      case 'F7':
        e.preventDefault();
        holdCurrentOrder();
        break;

      case 'F8':
        e.preventDefault();
        state.isReturnMode = !state.isReturnMode;
        if (state.isReturnMode) {
          dom.btnReturnMode.style.backgroundColor = 'var(--btn-dark)';
          dom.btnReturnMode.textContent = '[F8] ĐANG TRẢ HÀNG (-SL)';
          showToast('CHẾ ĐỘ TRẢ HÀNG: Mỗi món quét vào sẽ mang số lượng ÂM (-)', 'warning');
        } else {
          dom.btnReturnMode.style.backgroundColor = 'var(--btn-rose)';
          dom.btnReturnMode.textContent = '[F8] Trả hàng';
          showToast('Đã quay lại chế độ Bán hàng bình thường', 'info');
        }
        break;

      case 'F9':
        e.preventDefault();
        openDraftsModal();
        break;

      case 'F10':
        e.preventDefault();
        showToast('Tra cứu lịch sử hóa đơn bán lẻ [F10]', 'info');
        break;

      case 'F11':
        e.preventDefault();
        showToast(`Đang in lại hóa đơn [${state.currentSale.header.sct}] qua máy in bill...`, 'info');
        break;

      case 'F12':
        e.preventDefault();
        openModal(dom.modalLock);
        dom.inputUnlockPass.focus();
        break;

      case 'ArrowUp':
        if (document.activeElement !== dom.barcodeInput) {
          e.preventDefault();
          if (state.selectedItemIndex > 0) {
            state.selectedItemIndex--;
            renderCartHighlight();
          }
        }
        break;

      case 'ArrowDown':
        if (document.activeElement !== dom.barcodeInput) {
          e.preventDefault();
          if (state.selectedItemIndex < state.currentSale.items.length - 1) {
            state.selectedItemIndex++;
            renderCartHighlight();
          }
        }
        break;

      case 'Delete':
        if (document.activeElement !== dom.barcodeInput && state.selectedItemIndex >= 0) {
          e.preventDefault();
          removeCartItem(state.selectedItemIndex);
        }
        break;
    }
  }

  // =========================================================================
  // INITIALIZATION & EVENT LISTENERS
  // =========================================================================
  async function init() {
    // Clock
    setInterval(updateClock, 1000);
    updateClock();

    // Session info
    try {
      const sess = await api.getSessionInfo();
      if (sess) {
        state.session = sess;
        dom.sessionShift.textContent = sess.ca || 'Ca Sáng';
        dom.sessionCashier.textContent = 'Thu ngân: ' + (sess.thungan || 'ntd');
        dom.sessionStore.textContent = `Kho: ${sess.makho} (Quầy ${sess.quay})`;
      }
    } catch (e) {
      console.warn('Session load err', e);
    }

    // Categories
    try {
      state.categories = await api.getCategories();
      renderCategories();
    } catch (e) {
      console.warn('Categories err', e);
    }

    // Catalog products
    await loadCatalogProducts();

    // Start New Sale
    await startNewSale();

    // Event Listeners
    dom.barcodeInput.addEventListener('keydown', handleSearchInput);
    dom.barcodeInput.addEventListener('input', handleSearchInput);

    dom.catalogFilterInput.addEventListener('input', renderCatalogGrid);

    // Modal close buttons
    document.querySelectorAll('[data-close]').forEach(btn => {
      btn.addEventListener('click', () => {
        const modalId = btn.dataset.close;
        const modal = document.getElementById(modalId);
        if (modal) closeModal(modal);
      });
    });

    // Action bar buttons
    dom.btnChangeQty.addEventListener('click', openQtyModal);
    dom.btnDiscount.addEventListener('click', openDiscountModal);
    dom.btnDiscountTag.addEventListener('click', openDiscountModal);
    dom.btnHoldOrder.addEventListener('click', holdCurrentOrder);
    dom.btnPendingDrafts.addEventListener('click', openDraftsModal);
    dom.btnReturnMode.addEventListener('click', () => {
      const e = new KeyboardEvent('keydown', { key: 'F8' });
      handleGlobalKeyDown(e);
    });
    dom.btnCancelOrder.addEventListener('click', () => {
      if (confirm('Bạn có chắc chắn muốn hủy đơn hàng này không?')) {
        startNewSale();
      }
    });
    dom.btnCheckout.addEventListener('click', openPaymentModal);

    // Top action buttons
    dom.btnSearchOrder.addEventListener('click', () => {
      showToast('Tính năng tra cứu nhanh lịch sử đơn hàng', 'info');
    });
    dom.btnReprint.addEventListener('click', () => {
      showToast('Đang gửi lệnh in lại bill tới máy in nhiệt...', 'info');
    });
    dom.btnLock.addEventListener('click', () => {
      openModal(dom.modalLock);
      dom.inputUnlockPass.focus();
    });

    // Customer
    dom.btnCustomer.addEventListener('click', openCustomerModal);
    dom.inputSearchCustomer.addEventListener('input', searchCustomers);
    dom.btnChooseWalkin.addEventListener('click', () => {
      setCustomer({ madt: '', tendt: 'Khách lẻ', dienthoai: '', mavip: '' });
      closeModal(dom.modalCustomer);
    });

    // Payment modal inputs & quick cash
    dom.inputCashGiven.addEventListener('input', updatePaymentChange);
    dom.modalPayment.querySelectorAll('.payment-method-btn').forEach(btn => {
      btn.addEventListener('click', () => selectPaymentMethod(btn.dataset.method));
    });
    dom.modalPayment.querySelectorAll('.quick-cash-btn').forEach(btn => {
      btn.addEventListener('click', () => {
        const amt = btn.dataset.amt;
        if (amt === 'exact') {
          dom.inputCashGiven.value = Math.round(state.currentSale.header.thanhtoan).toLocaleString('vi-VN');
        } else {
          dom.inputCashGiven.value = Number(amt).toLocaleString('vi-VN');
        }
        updatePaymentChange();
      });
    });
    dom.btnSubmitPay.addEventListener('click', executePayment);

    // Qty Modal quick steps
    dom.modalQty.querySelectorAll('.quick-qty-btn').forEach(btn => {
      btn.addEventListener('click', () => {
        const step = parseFloat(btn.dataset.step) || 0;
        const cur = parseFloat(dom.inputNewQty.value) || 0;
        dom.inputNewQty.value = Math.max(1, cur + step);
      });
    });
    dom.btnSubmitQty.addEventListener('click', submitQtyModal);

    // Discount modal radio toggle
    dom.modalDiscount.querySelectorAll('input[name="discount-type"]').forEach(r => {
      r.addEventListener('change', () => {
        if (r.value === 'percent') {
          dom.discountInputLabel.textContent = 'Tỷ lệ giảm (%):';
          dom.inputDiscountVal.placeholder = 'Ví dụ: 10';
          dom.inputDiscountVal.max = '100';
        } else {
          dom.discountInputLabel.textContent = 'Số tiền giảm (VNĐ):';
          dom.inputDiscountVal.placeholder = 'Ví dụ: 50000';
          dom.inputDiscountVal.removeAttribute('max');
        }
      });
    });
    dom.btnSubmitDiscount.addEventListener('click', submitDiscountModal);

    // Unlock button
    dom.btnUnlock.addEventListener('click', () => {
      closeModal(dom.modalLock);
      dom.inputUnlockPass.value = '';
      showToast('Đã mở khóa màn hình thu ngân', 'success');
    });

    // Global Key Listener
    window.addEventListener('keydown', handleGlobalKeyDown);

    // Click outside search dropdown to close
    document.addEventListener('click', (e) => {
      if (!dom.barcodeInput.contains(e.target) && !dom.searchResultsDropdown.contains(e.target)) {
        closeSearchDropdown();
      }
    });

    console.log('✓ ShopHoa POS UI Initialized. Wails Mode:', state.isWails);
  }

  async function loadCatalogProducts() {
    try {
      state.catalogItems = await api.searchProducts('', state.selectedCategory);
      renderCatalogGrid();
    } catch (e) {
      console.warn('Catalog load err', e);
    }
  }

  // Start app when DOM is ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }

})();
