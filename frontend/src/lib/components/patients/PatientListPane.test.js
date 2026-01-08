import { render, screen, fireEvent } from "@testing-library/svelte";
import { describe, it, expect, vi } from "vitest";
import { createI18nMock } from "../../test-utils/i18nMock.js";

vi.mock("svelte-i18n", () => createI18nMock({
  "patients.noPatients": "No patients",
  "patients.noPatientsMatch": "No patients match",
  "patients.getStarted": "Get started",
  "patients.addPatient": "Add Patient",
  "patients.reminders": "Reminders",
  "common.edit": "Edit",
  "common.delete": "Delete",
  "common.searchPlaceholder": "Search...",
}));

import PatientListPane from "./PatientListPane.svelte";

describe("PatientListPane", () => {
  const mockPatients = [
    { id: "1", name: "John Doe", phone: "628123456789", email: "john@example.com", reminders: [{ id: "r1" }] },
    { id: "2", name: "Jane Smith", phone: "628987654321", email: "jane@example.com", reminders: [] }
  ];

  const defaultProps = {
    patients: mockPatients,
    selectedPatientId: null,
    searchQuery: "",
    onSelect: vi.fn(),
    onAddPatient: vi.fn(),
    onEditPatient: null,
    onDeletePatient: null
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("Rendering", () => {
    it("renders patient list with patients", () => {
      render(PatientListPane, defaultProps);
      expect(screen.getByText("John Doe")).toBeInTheDocument();
      expect(screen.getByText("Jane Smith")).toBeInTheDocument();
    });

    it("renders reminder count badge", () => {
      render(PatientListPane, defaultProps);
      expect(screen.getByText("1")).toBeInTheDocument();
    });

    it("renders add patient button", () => {
      render(PatientListPane, defaultProps);
      expect(screen.getByText("Add Patient")).toBeInTheDocument();
    });

    it("renders empty state when no patients", () => {
      render(PatientListPane, { ...defaultProps, patients: [] });
      expect(screen.getByText("No patients")).toBeInTheDocument();
    });
  });

  describe("Selection", () => {
    it("highlights selected patient", () => {
      render(PatientListPane, { ...defaultProps, selectedPatientId: "1" });
      const johnRow = screen.getByRole("option", { name: "John Doe" });
      expect(johnRow).toHaveAttribute("aria-selected", "true");
    });

    it("calls onSelect when patient is clicked", async () => {
      render(PatientListPane, defaultProps);
      const johnRow = screen.getByRole("option", { name: "John Doe" });
      await fireEvent.click(johnRow);
      expect(defaultProps.onSelect).toHaveBeenCalledWith("1");
    });
  });

  describe("Actions", () => {
    it("calls onAddPatient when add button is clicked", async () => {
      render(PatientListPane, defaultProps);
      const addButton = screen.getByText("Add Patient");
      await fireEvent.click(addButton);
      expect(defaultProps.onAddPatient).toHaveBeenCalled();
    });

    it("shows action buttons when patient is selected", () => {
      render(PatientListPane, {
        ...defaultProps,
        selectedPatientId: "1",
        onEditPatient: vi.fn(),
        onDeletePatient: vi.fn()
      });
      expect(screen.getByText("Edit")).toBeInTheDocument();
      expect(screen.getByText("Delete")).toBeInTheDocument();
    });

    it("calls onEditPatient when edit is clicked", async () => {
      const onEditPatient = vi.fn();
      render(PatientListPane, { ...defaultProps, selectedPatientId: "1", onEditPatient });
      const editButton = screen.getByText("Edit");
      await fireEvent.click(editButton);
      expect(onEditPatient).toHaveBeenCalledWith("1");
    });
  });

  describe("Accessibility", () => {
    it("has listbox role", () => {
      render(PatientListPane, defaultProps);
      expect(screen.getByRole("listbox")).toBeInTheDocument();
    });

    it("patient items have option role", () => {
      render(PatientListPane, defaultProps);
      expect(screen.getByRole("option", { name: "John Doe" })).toBeInTheDocument();
    });
  });
});
