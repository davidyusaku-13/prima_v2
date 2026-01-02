import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/svelte";
import CancelConfirmationModal from "./CancelConfirmationModal.svelte";
import { init, locale } from "svelte-i18n";

// Initialize i18n for tests
beforeEach(async () => {
  await init({
    fallbackLocale: "en",
    initialLocale: "en",
    messages: {
      en: {
        "reminder.cancel.title": "Cancel Reminder?",
        "reminder.cancel.confirm":
          "Are you sure you want to cancel this reminder?",
        "reminder.cancel.button": "Yes, Cancel",
        "reminder.cancel.cancelling": "Cancelling...",
        "reminder.cancel.scheduledTime": "Scheduled",
        "common.cancel": "Cancel",
      },
    },
  });
  locale.set("en");
});

describe("CancelConfirmationModal", () => {
  const mockReminder = {
    id: "reminder-1",
    title: "Test Reminder",
    description: "Test description",
    scheduled_at: "2026-01-15T10:00:00Z",
    delivery_status: "pending",
  };

  it("should not render when show is false", () => {
    const { container } = render(CancelConfirmationModal, {
      props: {
        show: false,
        reminder: null,
        onClose: vi.fn(),
        onConfirm: vi.fn(),
      },
    });

    expect(container.querySelector(".fixed")).toBeNull();
  });

  it("should render modal when show is true", async () => {
    render(CancelConfirmationModal, {
      props: {
        show: true,
        reminder: mockReminder,
        onClose: vi.fn(),
        onConfirm: vi.fn(),
      },
    });

    await waitFor(() => {
      expect(screen.getByText("Yes, Cancel")).toBeInTheDocument();
    });
    expect(screen.getByText("Test Reminder")).toBeInTheDocument();
    expect(screen.getByText("Test description")).toBeInTheDocument();
  });

  it("should display scheduled time when present", async () => {
    render(CancelConfirmationModal, {
      props: {
        show: true,
        reminder: mockReminder,
        onClose: vi.fn(),
        onConfirm: vi.fn(),
      },
    });

    await waitFor(() => {
      expect(screen.getByText(/Scheduled/i)).toBeInTheDocument();
    });
  });

  it("should call onClose when backdrop is clicked", async () => {
    const onClose = vi.fn();
    const { container } = render(CancelConfirmationModal, {
      props: {
        show: true,
        reminder: mockReminder,
        onClose,
        onConfirm: vi.fn(),
      },
    });

    const backdrop = container.querySelector(".bg-slate-900\\/50");
    await fireEvent.click(backdrop);

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it("should call onClose when Cancel button is clicked", async () => {
    const onClose = vi.fn();
    render(CancelConfirmationModal, {
      props: {
        show: true,
        reminder: mockReminder,
        onClose,
        onConfirm: vi.fn(),
      },
    });

    await waitFor(() => {
      expect(screen.getByText("Cancel")).toBeInTheDocument();
    });

    const cancelButton = screen.getByText("Cancel");
    await fireEvent.click(cancelButton);

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it("should call onConfirm when Yes, Cancel button is clicked", async () => {
    const onConfirm = vi.fn().mockResolvedValue(undefined);
    render(CancelConfirmationModal, {
      props: {
        show: true,
        reminder: mockReminder,
        onClose: vi.fn(),
        onConfirm,
      },
    });

    await waitFor(() => {
      expect(screen.getByText("Yes, Cancel")).toBeInTheDocument();
    });

    const confirmButton = screen.getByText("Yes, Cancel");
    await fireEvent.click(confirmButton);

    expect(onConfirm).toHaveBeenCalledTimes(1);
  });

  it("should show loading state while confirming", async () => {
    let resolveConfirm;
    const onConfirm = vi.fn(
      () =>
        new Promise((resolve) => {
          resolveConfirm = resolve;
        })
    );

    render(CancelConfirmationModal, {
      props: {
        show: true,
        reminder: mockReminder,
        onClose: vi.fn(),
        onConfirm,
      },
    });

    await waitFor(() => {
      expect(screen.getByText("Yes, Cancel")).toBeInTheDocument();
    });

    const confirmButton = screen.getByText("Yes, Cancel");
    await fireEvent.click(confirmButton);

    // Should show loading text
    await waitFor(() => {
      expect(screen.getByText("Cancelling...")).toBeInTheDocument();
    });

    // Buttons should be disabled
    expect(confirmButton).toBeDisabled();

    // Resolve the promise
    resolveConfirm();
  });

  it("should close on Escape key press", async () => {
    const onClose = vi.fn();
    render(CancelConfirmationModal, {
      props: {
        show: true,
        reminder: mockReminder,
        onClose,
        onConfirm: vi.fn(),
      },
    });

    await waitFor(() => {
      expect(screen.getByText("Yes, Cancel")).toBeInTheDocument();
    });

    // Simulate ESC key press on window
    const escapeEvent = new KeyboardEvent("keydown", { key: "Escape" });
    window.dispatchEvent(escapeEvent);

    await waitFor(() => {
      expect(onClose).toHaveBeenCalledTimes(1);
    });
  });

  it("should not close on Escape when loading", async () => {
    let resolveConfirm;
    const onConfirm = vi.fn(
      () =>
        new Promise((resolve) => {
          resolveConfirm = resolve;
        })
    );
    const onClose = vi.fn();

    render(CancelConfirmationModal, {
      props: {
        show: true,
        reminder: mockReminder,
        onClose,
        onConfirm,
      },
    });

    await waitFor(() => {
      expect(screen.getByText("Yes, Cancel")).toBeInTheDocument();
    });

    // Click confirm to start loading
    const confirmButton = screen.getByText("Yes, Cancel");
    await fireEvent.click(confirmButton);

    await waitFor(() => {
      expect(screen.getByText("Cancelling...")).toBeInTheDocument();
    });

    // Try to close with ESC while loading
    const escapeEvent = new KeyboardEvent("keydown", { key: "Escape" });
    window.dispatchEvent(escapeEvent);

    // onClose should NOT be called
    expect(onClose).not.toHaveBeenCalled();

    // Cleanup
    resolveConfirm();
  });

  it("should handle reminder without scheduled_at but with due_date", async () => {
    const reminderWithDueDate = {
      ...mockReminder,
      scheduled_at: undefined,
      due_date: "2026-01-20T14:00:00Z",
    };

    render(CancelConfirmationModal, {
      props: {
        show: true,
        reminder: reminderWithDueDate,
        onClose: vi.fn(),
        onConfirm: vi.fn(),
      },
    });

    await waitFor(() => {
      expect(screen.getByText(/Scheduled/i)).toBeInTheDocument();
    });
  });

  it("should display message when description is not present", async () => {
    const reminderWithMessage = {
      ...mockReminder,
      description: undefined,
      message: "Custom message text",
    };

    render(CancelConfirmationModal, {
      props: {
        show: true,
        reminder: reminderWithMessage,
        onClose: vi.fn(),
        onConfirm: vi.fn(),
      },
    });

    await waitFor(() => {
      expect(screen.getByText("Custom message text")).toBeInTheDocument();
    });
  });
});
