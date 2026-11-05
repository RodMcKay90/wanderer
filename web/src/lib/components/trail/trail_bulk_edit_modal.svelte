<script module lang="ts">
    export type TrailBulkEditChanges = {
        category?: {
            category: string;
            subcategory: string;
        };
        difficulty?: "easy" | "moderate" | "difficult";
        route_type?: "loop" | "traversal";
        ride_profile?: string;
        starts_home?: boolean;
    };
</script>

<script lang="ts">
    import Modal from "#lib/components/base/modal.svelte";
    import Select from "#lib/components/base/select.svelte";
    import CategoryPicker, {
        type CategoryPickerSelection,
    } from "./category_picker.svelte";
    import { _ } from "svelte-i18n";
    import {
        RIDE_PROFILES,
        ROUTE_TYPES,
        type RouteType,
    } from "#lib/util/ride_attributes.ts";

    interface Props {
        selectedCount?: number;
        initialCategorySelection?: CategoryPickerSelection;
        initialDifficulty?: TrailBulkEditChanges["difficulty"];
        onapply?: (changes: TrailBulkEditChanges) => Promise<void> | void;
    }

    let {
        selectedCount = 0,
        initialCategorySelection,
        initialDifficulty,
        onapply,
    }: Props = $props();

    const uid = $props.id();
    const modalId = `trail-bulk-edit-modal-${uid}`;

    let modal: Modal;
    let applyCategory = $state(false);
    let applyRouteType = $state(false);
    let applyRideProfile = $state(false);
    let applyStartsHome = $state(false);
    let categorySelection: CategoryPickerSelection | undefined = $state();
    let categoryValue = $state("");
    let routeTypeValue: RouteType = $state("loop");
    let rideProfileValue: string = $state(RIDE_PROFILES[0].value);
    let startsHomeValue = $state(true);
    let loading = $state(false);

    let canApply = $derived(
        (applyCategory && categorySelection !== undefined) ||
            applyRouteType ||
            applyRideProfile ||
            applyStartsHome,
    );

    export function openModal() {
        applyCategory = false;
        applyRouteType = false;
        applyRideProfile = false;
        applyStartsHome = false;
        categorySelection = initialCategorySelection;
        categoryValue = categoryValueFromSelection(initialCategorySelection);
        modal.openModal();
    }

    function categoryValueFromSelection(selection?: CategoryPickerSelection) {
        if (!selection) {
            return "";
        }

        return selection.subcategory
            ? `subcategory:${selection.subcategory}`
            : `category:${selection.category}`;
    }

    function handleCategoryChange(selection: CategoryPickerSelection) {
        categorySelection = selection;
        categoryValue = selection.subcategory
            ? `subcategory:${selection.subcategory}`
            : `category:${selection.category}`;
    }

    async function apply(closeModal: () => void) {
        if (!canApply) {
            return;
        }

        const changes: TrailBulkEditChanges = {};
        if (applyCategory) {
            if (!categorySelection) {
                return;
            }
            changes.category = categorySelection;
        }
        if (applyRouteType) {
            changes.route_type = routeTypeValue;
        }
        if (applyRideProfile) {
            changes.ride_profile = rideProfileValue;
        }
        if (applyStartsHome) {
            changes.starts_home = startsHomeValue;
        }

        loading = true;
        try {
            await onapply?.(changes);
            closeModal();
        } finally {
            loading = false;
        }
    }
</script>

<Modal
    id={modalId}
    title={$_("adjust")}
    size="md:min-w-sm"
    bind:this={modal}
>
    {#snippet content()}
        <p class="text-sm text-gray-500">
            {$_("bulk-edit-selected-trails", {
                values: { n: selectedCount },
            })}
        </p>

        <div class="space-y-4">
            <label class="flex items-center gap-2">
                <input
                    type="checkbox"
                    class="h-4 w-4 accent-primary"
                    bind:checked={applyCategory}
                />
                <span class="font-medium">{$_("category")}</span>
            </label>
            {#if applyCategory}
                <CategoryPicker
                    value={categoryValue}
                    label=""
                    loadData
                    fixedDropdown
                    onchange={handleCategoryChange}
                ></CategoryPicker>
            {/if}

            <label class="flex items-center gap-2">
                <input
                    type="checkbox"
                    class="h-4 w-4 accent-primary"
                    bind:checked={applyRouteType}
                />
                <span class="font-medium">Streckentyp</span>
            </label>
            {#if applyRouteType}
                <Select
                    value={routeTypeValue}
                    onchange={(value) => (routeTypeValue = value)}
                    items={ROUTE_TYPES.map((r) => ({
                        text: r.text,
                        value: r.value,
                    }))}
                ></Select>
            {/if}

            <label class="flex items-center gap-2">
                <input
                    type="checkbox"
                    class="h-4 w-4 accent-primary"
                    bind:checked={applyRideProfile}
                />
                <span class="font-medium">Profil</span>
            </label>
            {#if applyRideProfile}
                <Select
                    value={rideProfileValue}
                    onchange={(value) => (rideProfileValue = value)}
                    items={RIDE_PROFILES}
                ></Select>
            {/if}

            <label class="flex items-center gap-2">
                <input
                    type="checkbox"
                    class="h-4 w-4 accent-primary"
                    bind:checked={applyStartsHome}
                />
                <span class="font-medium">Start zuhause</span>
            </label>
            {#if applyStartsHome}
                <Select
                    value={startsHomeValue}
                    onchange={(value) =>
                        (startsHomeValue = value === true || value === "true")}
                    items={[
                        { text: "Ja", value: true },
                        { text: "Nein", value: false },
                    ]}
                ></Select>
            {/if}
        </div>
    {/snippet}

    {#snippet footer({ closeModal })}
        <div class="flex items-center gap-4">
            <button class="btn-secondary" type="button" onclick={closeModal}>
                {$_("cancel")}
            </button>
            <button
                class="btn-primary"
                class:btn-disabled={!canApply || loading}
                disabled={!canApply || loading}
                type="button"
                onclick={() => apply(closeModal)}
            >
                {$_("apply")}
            </button>
        </div>
    {/snippet}
</Modal>
