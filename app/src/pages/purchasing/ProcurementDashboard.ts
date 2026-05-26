import { LitElement, html, nothing } from 'lit';
import { customElement, state } from 'lit/decorators.js';
import { icon } from '../../lib/icons.ts';
import { ArrowLeft, RefreshCw, CheckCircle, XCircle, Package, TrendingUp, Clock, ShoppingCart } from 'lucide';
import { ToastService } from '../../lib/toast-service.ts';
import { router } from '../../lib/router.ts';
import type { DashboardSummary, VendorDraftGroup, ReplenishmentDraft } from '../../types/replenishment.ts';

/**
 * Procurement Dashboard — Human-in-the-Loop Review Panel
 *
 * Displays draft POs grouped by manufacturer/vendor as expandable cards.
 * Procurement manager reviews, edits quantities, and approves/rejects
 * suggested purchase orders before they are sent to suppliers.
 *
 * Route: /purchasing/procurement-dashboard
 *
 * TODO: Implement full render and interaction logic
 */
@customElement('gable-procurement-dashboard')
export class GableProcurementDashboard extends LitElement {
    createRenderRoot() { return this; }

    @state() private dashboard: DashboardSummary | null = null;
    @state() private loading = true;
    @state() private generating = false;
    @state() private expandedVendor: string | null = null;

    connectedCallback() {
        super.connectedCallback();
        this._loadDashboard();
    }

    private async _loadDashboard() {
        this.loading = true;
        try {
            // TODO: Call PurchaseOrderService.getProcurementDashboard()
            // this.dashboard = await PurchaseOrderService.getProcurementDashboard();
            this.dashboard = null; // Stub
        } catch (err) {
            console.error('Failed to load procurement dashboard:', err);
            ToastService.show('Failed to load procurement dashboard', 'error');
        } finally {
            this.loading = false;
        }
    }

    private async _handleGenerate() {
        this.generating = true;
        try {
            // TODO: Call PurchaseOrderService.generateProcurementDrafts()
            ToastService.show('Draft POs generated successfully', 'success');
            await this._loadDashboard();
        } catch (err) {
            console.error('Failed to generate drafts:', err);
            ToastService.show('Failed to generate draft POs', 'error');
        } finally {
            this.generating = false;
        }
    }

    private async _handleApprove(draft: ReplenishmentDraft) {
        try {
            // TODO: Call PurchaseOrderService.approveProcurementDraft(draft.id)
            ToastService.show(`Draft PO approved and sent to ${draft.vendor_name}`, 'success');
            await this._loadDashboard();
        } catch (err) {
            console.error('Failed to approve draft:', err);
            ToastService.show('Failed to approve draft', 'error');
        }
    }

    private async _handleReject(_draft: ReplenishmentDraft) {
        // TODO: Show rejection notes modal before calling reject endpoint
        try {
            // TODO: Call PurchaseOrderService.rejectProcurementDraft(draft.id, notes)
            ToastService.show(`Draft PO rejected`, 'info');
            await this._loadDashboard();
        } catch (err) {
            console.error('Failed to reject draft:', err);
            ToastService.show('Failed to reject draft', 'error');
        }
    }

    render() {
        return html`
            <div>
                <!-- Header -->
                <div class="flex items-center gap-4 mb-6">
                    <button
                        @click=${() => router.navigate('/purchasing')}
                        class="p-2 rounded-full bg-white/5 hover:bg-white/10 text-zinc-400 transition-colors"
                    >
                        ${icon(ArrowLeft, 20, 'w-5 h-5')}
                    </button>
                    <div class="flex-1">
                        <h1 class="text-2xl font-bold text-white">Procurement Dashboard</h1>
                        <p class="text-sm text-zinc-400">Review and approve suggested purchase orders before sending to suppliers</p>
                    </div>
                    <button
                        @click=${() => this._handleGenerate()}
                        ?disabled=${this.generating}
                        class="flex items-center gap-2 px-4 py-2 bg-[#00FFA3] text-black text-sm font-semibold rounded-lg hover:bg-[#00FFA3]/90 transition-colors disabled:opacity-50"
                    >
                        ${icon(RefreshCw, 16, `w-4 h-4 ${this.generating ? 'animate-spin' : ''}`)}
                        Generate Drafts
                    </button>
                </div>

                <!-- Summary KPIs -->
                ${this.dashboard ? html`
                    <div class="grid grid-cols-2 md:grid-cols-3 gap-4 mb-6">
                        <div class="bg-[#161821]/60 backdrop-blur-sm border border-white/10 rounded-2xl p-4 text-center">
                            ${icon(Package, 20, 'w-5 h-5 mx-auto mb-1 text-zinc-400')}
                            <div class="text-2xl font-bold text-white font-mono">${this.dashboard.pending_count}</div>
                            <div class="text-xs text-zinc-500">Pending Drafts</div>
                        </div>
                        <div class="bg-[#161821]/60 backdrop-blur-sm border border-white/10 rounded-2xl p-4 text-center">
                            ${icon(ShoppingCart, 20, 'w-5 h-5 mx-auto mb-1 text-[#00FFA3]')}
                            <div class="text-2xl font-bold text-[#00FFA3] font-mono">
                                $${this.dashboard.total_est_cost.toLocaleString('en-US', { minimumFractionDigits: 0 })}
                            </div>
                            <div class="text-xs text-zinc-500">Est. Total Cost</div>
                        </div>
                        <div class="bg-[#161821]/60 backdrop-blur-sm border border-white/10 rounded-2xl p-4 text-center">
                            ${icon(TrendingUp, 20, 'w-5 h-5 mx-auto mb-1 text-blue-400')}
                            <div class="text-2xl font-bold text-blue-400 font-mono">${this.dashboard.vendor_groups.length}</div>
                            <div class="text-xs text-zinc-500">Vendors</div>
                        </div>
                    </div>
                ` : nothing}

                <!-- Loading State -->
                ${this.loading ? html`
                    <div class="bg-[#161821]/60 backdrop-blur-sm border border-white/10 rounded-2xl">
                        <div class="text-center text-zinc-500 py-16">
                            ${icon(RefreshCw, 32, 'w-8 h-8 mx-auto mb-3 animate-spin text-zinc-600')}
                            Loading procurement dashboard...
                        </div>
                    </div>
                ` : nothing}

                <!-- Empty State -->
                ${!this.loading && (!this.dashboard || this.dashboard.vendor_groups.length === 0) ? html`
                    <div class="bg-[#161821]/60 backdrop-blur-sm border border-white/10 rounded-2xl">
                        <div class="text-center py-16">
                            ${icon(Package, 48, 'w-12 h-12 mx-auto mb-4 text-zinc-700')}
                            <h3 class="text-lg font-medium text-white mb-2">No pending drafts</h3>
                            <p class="text-sm text-zinc-500 mb-6">Click "Generate Drafts" to create suggested purchase orders<br/>based on current inventory and sales velocity.</p>
                        </div>
                    </div>
                ` : nothing}

                <!-- Vendor Groups -->
                ${!this.loading && this.dashboard?.vendor_groups.map(group => this._renderVendorGroup(group))}
            </div>
        `;
    }

    private _renderVendorGroup(group: VendorDraftGroup) {
        const isExpanded = this.expandedVendor === group.vendor_id;
        return html`
            <div class="bg-[#161821]/60 backdrop-blur-sm border border-white/10 rounded-2xl mb-4 overflow-hidden">
                <!-- Vendor Header (clickable to expand) -->
                <button
                    @click=${() => { this.expandedVendor = isExpanded ? null : group.vendor_id; }}
                    class="w-full flex items-center justify-between p-4 hover:bg-white/[0.02] transition-colors text-left"
                >
                    <div class="flex items-center gap-3">
                        <div class="text-lg font-semibold text-white">${group.vendor_name}</div>
                        <span class="text-xs text-zinc-500 flex items-center gap-1">
                            ${icon(Clock, 12, 'w-3 h-3')} ${group.lead_time_days.toFixed(0)}d lead time
                        </span>
                        <span class="text-xs px-2 py-0.5 rounded bg-white/5 text-zinc-400">
                            ${group.drafts.length} draft${group.drafts.length !== 1 ? 's' : ''}
                        </span>
                    </div>
                    <div class="text-sm font-mono text-[#00FFA3] font-medium">
                        $${group.total_cost.toLocaleString('en-US', { minimumFractionDigits: 2 })}
                    </div>
                </button>

                <!-- Expanded Draft Details -->
                ${isExpanded ? html`
                    <div class="border-t border-white/5">
                        ${group.drafts.map(draft => this._renderDraft(draft))}
                    </div>
                ` : nothing}
            </div>
        `;
    }

    private _renderDraft(draft: ReplenishmentDraft) {
        // TODO: Render individual draft with editable line items, confidence indicator,
        //       and Approve/Reject action buttons
        return html`
            <div class="p-4 border-b border-white/5 last:border-b-0">
                <div class="flex items-center justify-between mb-3">
                    <div class="flex items-center gap-3">
                        <span class="text-sm font-medium text-white">${draft.total_lines} items</span>
                        <span class="text-xs text-zinc-500">
                            Confidence: <span class="font-mono ${draft.confidence >= 80 ? 'text-emerald-400' : draft.confidence >= 50 ? 'text-amber-400' : 'text-rose-400'}">${draft.confidence.toFixed(0)}%</span>
                        </span>
                    </div>
                    <div class="flex items-center gap-2">
                        <button
                            @click=${() => this._handleApprove(draft)}
                            class="flex items-center gap-1.5 px-3 py-1.5 bg-[#00FFA3] text-black text-xs font-semibold rounded-lg hover:bg-[#00FFA3]/90 transition-colors"
                        >
                            ${icon(CheckCircle, 14, 'w-3.5 h-3.5')} Approve & Send
                        </button>
                        <button
                            @click=${() => this._handleReject(draft)}
                            class="flex items-center gap-1.5 px-3 py-1.5 bg-white/5 border border-white/10 text-zinc-400 text-xs font-medium rounded-lg hover:bg-white/10 transition-colors"
                        >
                            ${icon(XCircle, 14, 'w-3.5 h-3.5')} Reject
                        </button>
                    </div>
                </div>

                <!-- TODO: Render editable PO lines table here -->
                <div class="text-xs text-zinc-600 italic py-4 text-center">
                    Draft PO line items will render here with editable quantities
                </div>
            </div>
        `;
    }
}
