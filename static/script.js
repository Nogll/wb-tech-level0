document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('searchForm');
    const uidInput = document.getElementById('uidInput');
    const error = document.getElementById('error');
    const orderInfo = document.getElementById('orderInfo');
    const itemsContainer = document.getElementById('itemsContainer');

    form.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const uid = uidInput.value.trim();
        if (!uid) {
            showError('Пожалуйста, введите UID заказа');
            return;
        }

        await fetchOrder(uid);
    });

    async function fetchOrder(uid) {
        try {
            hideError();
            hideOrderInfo();

            // Замените URL на ваш реальный endpoint
            const response = await fetch(`/api/v1/orders?uid=${encodeURIComponent(uid)}`);
            
            if (!response.ok) {
                throw new Error(`Ошибка HTTP: ${response.status}`);
            }

            const order = await response.json();
            displayOrder(order);
            
        } catch (err) {
            showError(`Ошибка при получении данных: ${err.message}`);
        } finally {
        }
    }

    function displayOrder(order) {
        // Заполняем основную информацию о заказе
        document.getElementById('orderUID').textContent = order.order_uid || '—';
        document.getElementById('trackNumber').textContent = order.track_number || '—';
        document.getElementById('entry').textContent = order.entry || '—';
        document.getElementById('locale').textContent = order.locale || '—';
        document.getElementById('internalSignature').textContent = order.internal_signature || '—';
        document.getElementById('customerID').textContent = order.customer_id || '—';
        document.getElementById('deliveryService').textContent = order.delivery_service || '—';
        document.getElementById('shardKey').textContent = order.shardkey || '—';
        document.getElementById('smID').textContent = order.sm_id || '—';
        document.getElementById('dateCreated').textContent = order.date_created || '—';
        document.getElementById('oofShard').textContent = order.oof_shard || '—';

        // Заполняем информацию о доставке
        if (order.delivery) {
            document.getElementById('deliveryName').textContent = order.delivery.name || '—';
            document.getElementById('deliveryPhone').textContent = order.delivery.phone || '—';
            document.getElementById('deliveryZip').textContent = order.delivery.zip || '—';
            document.getElementById('deliveryCity').textContent = order.delivery.city || '—';
            document.getElementById('deliveryAddress').textContent = order.delivery.address || '—';
            document.getElementById('deliveryRegion').textContent = order.delivery.region || '—';
            document.getElementById('deliveryEmail').textContent = order.delivery.email || '—';
        }

        // Заполняем информацию о платеже
        if (order.payment) {
            document.getElementById('paymentTransaction').textContent = order.payment.transaction || '—';
            document.getElementById('paymentRequestID').textContent = order.payment.request_id || '—';
            document.getElementById('paymentCurrency').textContent = order.payment.currency || '—';
            document.getElementById('paymentProvider').textContent = order.payment.provider || '—';
            document.getElementById('paymentAmount').textContent = order.payment.amount || '—';
            document.getElementById('paymentDT').textContent = formatTimestamp(order.payment.payment_dt) || '—';
            document.getElementById('paymentBank').textContent = order.payment.bank || '—';
            document.getElementById('paymentDeliveryCost').textContent = order.payment.delivery_cost || '—';
            document.getElementById('paymentGoodsTotal').textContent = order.payment.goods_total || '—';
            document.getElementById('paymentCustomFee').textContent = order.payment.custom_fee || '0';
        }

        // Заполняем информацию о товарах
        displayItems(order.items || []);

        showOrderInfo();
    }

    function displayItems(items) {
        itemsContainer.innerHTML = '';
        
        if (items.length === 0) {
            itemsContainer.innerHTML = '<p>Товары не найдены</p>';
            return;
        }

        items.forEach((item, index) => {
            const itemElement = document.createElement('div');
            itemElement.className = 'item';
            itemElement.innerHTML = `
                <h4>Товар ${index + 1}</h4>
                <table>
                    <tr><td>ChrtID</td><td>${item.chrt_id || '—'}</td></tr>
                    <tr><td>Трэк номер</td><td>${item.track_number || '—'}</td></tr>
                    <tr><td>Price</td><td>${item.price || '—'}</td></tr>
                    <tr><td>RID</td><td>${item.rid || '—'}</td></tr>
                    <tr><td>Name</td><td>${item.name || '—'}</td></tr>
                    <tr><td>Sale</td><td>${item.sale || '—'}</td></tr>
                    <tr><td>Size</td><td>${item.size || '—'}</td></tr>
                    <tr><td>Total Price</td><td>${item.total_price || '—'}</td></tr>
                    <tr><td>NmID</td><td>${item.nm_id || '—'}</td></tr>
                    <tr><td>Brand</td><td>${item.brand || '—'}</td></tr>
                    <tr><td>Status</td><td>${item.status || '—'}</td></tr>
                </table>
            `;
            itemsContainer.appendChild(itemElement);
        });
    }

    function formatTimestamp(timestamp) {
        if (!timestamp) return '—';
        return new Date(timestamp * 1000).toLocaleString();
    }

    function showError(message) {
        error.classList.remove('hidden');
        document.getElementById('errorMessage').textContent = message;
    }

    function hideError() {
        error.classList.add('hidden');
    }

    function showOrderInfo() {
        orderInfo.classList.remove('hidden');
    }

    function hideOrderInfo() {
        orderInfo.classList.add('hidden');
    }
});