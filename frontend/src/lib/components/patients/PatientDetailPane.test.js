import { render, screen, fireEvent } from "@testing-library/svelte";
import { describe, it, expect, vi } from "vitest";
import { createI18nMock } from "../../test-utils/i18nMock.js";

// Mock DeliveryStatusBadge component
vi.mock("../delivery/DeliveryStatusBadge.svelte", () => ({
  default: vi.fn(() => null)
}));

// Mock ReminderHistoryView component
vi.mock("../../views/patients/ReminderHistoryView.svelte", () => ({
  default: vi.fn(() => null)
}));

// Mock toast store
vi.mock("../../stores/toast.svelte.js", () => ({
  toastStore: {
    add: vi.fn()
  }
}));

vi.mock("svelte-i18n", () => createI18nMock({
  "patients.selectPatient": "Select patient",
  "patients.selectPatientDescription": "Choose from list",
  "patients.info": "Info",
  "patients.reminders": "Reminders",
  "patients.reminderHistory": "History",
  "patients.patientName": "Name",
  "patients.phone": "Phone",
  "patients.email": "Email",
  "patients.notes": "Notes",
  "patients.createdAt": "Created",
  "patients.noReminders": "No reminders",
  "patients.noRemindersDescription": "Create reminder",
  "patients.addReminder": "Add Reminder",
  "reminder.retry": "Retry",
  "reminder.send.title": "Send",
  "common.edit": "Edit",
  "common.delete": "Delete",
  "patients.invalidPhoneError": "Invalid phone number"
}));

import PatientDetailPane from "./PatientDetailPane.svelte";

describe("PatientDetailPane", () => {
  const mockPatient = {
    id: "1",
    name: "John Doe",
    phone: "628123456789",
    email: "john@example.com",
    notes: "Test notes",
    createdAt: "2024-01-15T10:00:00Z",
    reminders: [
      { id: "r1", title: "Reminder 1", priority: "high", delivery_status: "pending" },
      { id: "r2", title: "Reminder 2", priority: "low", delivery_status: "sent" }
    ]
  };

  const defaultProps = {
    patient: mockPatient,
    deliveryStatuses: {},
    failedReminders: [],
    activeTab: "info",
    onTabChange: vi.fn()
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("Empty State", () => {
    it("shows empty state when no patient is selected", () => {
      render(PatientDetailPane, { patient: null });
      expect(screen.getByText("Select patient")).toBeInTheDocument();
    });
  });

  describe("Patient Header", () => {
    it("renders patient name in header", () => {
      render(PatientDetailPane, defaultProps);
      const nameElement = screen.getByRole("heading", { name: "John Doe" });
      expect(nameElement).toBeInTheDocument();
    });

    it("renders patient avatar with initial", () => {
      render(PatientDetailPane, defaultProps);
      expect(screen.getByText("J")).toBeInTheDocument();
    });

    it("renders WhatsApp button when phone exists", () => {
      render(PatientDetailPane, defaultProps);
      expect(screen.getByLabelText(/whatsapp/i)).toBeInTheDocument();
    });
  });

  describe("Tabs", () => {
    it("renders info tab", () => {
      render(PatientDetailPane, defaultProps);
      expect(screen.getByRole("tab", { name: /info/i })).toBeInTheDocument();
    });

    it("renders reminders tab", () => {
      render(PatientDetailPane, defaultProps);
      expect(screen.getByRole("tab", { name: /reminders/i })).toBeInTheDocument();
    });

    it("renders history tab", () => {
      render(PatientDetailPane, defaultProps);
      expect(screen.getByRole("tab", { name: /history/i })).toBeInTheDocument();
    });

    it("calls onTabChange when tab is clicked", async () => {
      render(PatientDetailPane, defaultProps);
      const remindersTab = screen.getByRole("tab", { name: /reminders/i });
      await fireEvent.click(remindersTab);
      expect(defaultProps.onTabChange).toHaveBeenCalledWith("reminders");
    });

    it("marks active tab as selected", () => {
      render(PatientDetailPane, { ...defaultProps, activeTab: "reminders" });
      const remindersTab = screen.getByRole("tab", { name: /reminders/i });
      expect(remindersTab).toHaveAttribute("aria-selected", "true");
    });
  });

  describe("Info Tab", () => {
    it("shows patient details in info tab", () => {
      render(PatientDetailPane, defaultProps);
      // Find the info panel specifically
      const infoPanel = document.getElementById("panel-info");
      expect(infoPanel).toBeInTheDocument();
      expect(infoPanel).toHaveTextContent("John Doe");
      expect(infoPanel).toHaveTextContent("628123456789");
    });
  });

  describe("Reminders Tab", () => {
    it("renders reminder list", () => {
      render(PatientDetailPane, { ...defaultProps, activeTab: "reminders" });
      expect(screen.getByText("Reminder 1")).toBeInTheDocument();
      expect(screen.getByText("Reminder 2")).toBeInTheDocument();
    });

    it("shows reminder priority badges", () => {
      render(PatientDetailPane, { ...defaultProps, activeTab: "reminders" });
      expect(screen.getByText("high")).toBeInTheDocument();
      expect(screen.getByText("low")).toBeInTheDocument();
    });
  });

  describe("Accessibility", () => {
    it("has tablist role", () => {
      render(PatientDetailPane, defaultProps);
      expect(screen.getByRole("tablist")).toBeInTheDocument();
    });

    it("tabs have tab role", () => {
      render(PatientDetailPane, defaultProps);
      expect(screen.getByRole("tab", { name: /info/i })).toHaveAttribute("role", "tab");
    });

    it("tab has aria-controls referencing panel", () => {
      render(PatientDetailPane, defaultProps);
      const infoTab = screen.getByRole("tab", { name: /info/i });
      expect(infoTab).toHaveAttribute("aria-controls", "panel-info");
    });
  });

  describe("WhatsApp Button", () => {
    it("opens WhatsApp with formatted phone", () => {
      const openStub = vi.spyOn(window, "open").mockImplementation(() => {});
      render(PatientDetailPane, defaultProps);
      const whatsappButton = screen.getByLabelText(/whatsapp/i);
      fireEvent.click(whatsappButton);
      expect(openStub).toHaveBeenCalledWith("https://wa.me/628123456789", "_blank");
      openStub.mockRestore();
    });

    it("shows error toast when phone is invalid", async () => {
      const { toastStore } = await import("../../stores/toast.svelte.js");
      const toastAddSpy = vi.spyOn(toastStore, "add");
      const mockPatientNoPhone = { ...mockPatient, phone: "" };
      render(PatientDetailPane, { ...defaultProps, patient: mockPatientNoPhone });
      const whatsappButton = screen.getByLabelText(/whatsapp|phone/i);
      await fireEvent.click(whatsappButton);
      expect(toastAddSpy).toHaveBeenCalledWith(expect.any(String), { type: "error" });
    });
  });

  describe("Tab Reset (AC26)", () => {
    it("resets tab to info when patient changes", async () => {
      const onTabChange = vi.fn();
      // Render with reminders tab active
      const { rerender } = render(PatientDetailPane, {
        ...defaultProps,
        activeTab: "reminders",
        onTabChange
      });

      // Switch to a different patient (using a new object reference)
      const newPatient = { ...mockPatient, id: "2", name: "Jane Doe" };
      await rerender({
        patient: newPatient,
        activeTab: "reminders", // Parent still thinks it's on reminders
        onTabChange
      });

      // Tab should reset to info because patient changed
      expect(onTabChange).toHaveBeenCalledWith("info");
    });
  });
});
