// File: src/components/auth/profile-step.tsx

"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Loader2 } from "lucide-react";
import { useState } from "react";

interface ProfileStepProps {
  onSubmit: (data: {
    firstName: string;
    lastName: string;
    phone?: string;
  }) => Promise<void> | void;
  isLoading: boolean;
  title?: string;
  buttonText?: string;
  showPhone?: boolean;
}

export function ProfileStep({
  onSubmit,
  isLoading,
  title = "Complete your profile",
  buttonText = "Create Account",
  showPhone = true,
}: ProfileStepProps) {
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [phone, setPhone] = useState("");
  const [errors, setErrors] = useState({
    firstName: "",
    lastName: "",
    phone: "",
  });

  const validateForm = () => {
    const newErrors = {
      firstName: "",
      lastName: "",
      phone: "",
    };

    // Validate first name
    if (!firstName.trim()) {
      newErrors.firstName = "First name is required";
    } else if (firstName.trim().length < 2) {
      newErrors.firstName = "First name must be at least 2 characters";
    } else if (!/^[a-zA-Z\s\-']+$/.test(firstName.trim())) {
      newErrors.firstName =
        "First name can only contain letters, spaces, hyphens, and apostrophes";
    }

    // Validate last name
    if (!lastName.trim()) {
      newErrors.lastName = "Last name is required";
    } else if (lastName.trim().length < 2) {
      newErrors.lastName = "Last name must be at least 2 characters";
    } else if (!/^[a-zA-Z\s\-']+$/.test(lastName.trim())) {
      newErrors.lastName =
        "Last name can only contain letters, spaces, hyphens, and apostrophes";
    }

    // Validate phone (optional)
    if (phone.trim() && showPhone) {
      const cleanPhone = phone.trim();
      if (!cleanPhone.startsWith("+")) {
        newErrors.phone = "Phone number must include country code (e.g., +91)";
      } else {
        const digits = cleanPhone.replace(/[^\d]/g, "");
        if (digits.length < 10 || digits.length > 15) {
          newErrors.phone = "Phone number must be between 10 and 15 digits";
        }
      }
    }

    setErrors(newErrors);
    return !newErrors.firstName && !newErrors.lastName && !newErrors.phone;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    const profileData = {
      firstName: firstName.trim(),
      lastName: lastName.trim(),
      phone: phone.trim() || undefined,
    };

    await onSubmit(profileData);
  };

  const handleFirstNameChange = (value: string) => {
    setFirstName(value);
    if (errors.firstName) {
      setErrors({ ...errors, firstName: "" });
    }
  };

  const handleLastNameChange = (value: string) => {
    setLastName(value);
    if (errors.lastName) {
      setErrors({ ...errors, lastName: "" });
    }
  };

  const handlePhoneChange = (value: string) => {
    setPhone(value);
    if (errors.phone) {
      setErrors({ ...errors, phone: "" });
    }
  };

  const isFormValid =
    firstName.trim() &&
    lastName.trim() &&
    !errors.firstName &&
    !errors.lastName &&
    !errors.phone;

  return (
    <div className="space-y-6">
      <div className="text-center">
        <p className="text-sm text-gray-600">
          Just a few more details to complete your account setup
        </p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label htmlFor="firstName">
              First Name <span className="text-red-500">*</span>
            </Label>
            <Input
              id="firstName"
              type="text"
              placeholder="John"
              value={firstName}
              onChange={(e) => handleFirstNameChange(e.target.value)}
              disabled={isLoading}
              required
            />
            {errors.firstName && (
              <p className="text-sm text-red-600">{errors.firstName}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="lastName">
              Last Name <span className="text-red-500">*</span>
            </Label>
            <Input
              id="lastName"
              type="text"
              placeholder="Doe"
              value={lastName}
              onChange={(e) => handleLastNameChange(e.target.value)}
              disabled={isLoading}
              required
            />
            {errors.lastName && (
              <p className="text-sm text-red-600">{errors.lastName}</p>
            )}
          </div>
        </div>

        {showPhone && (
          <div className="space-y-2">
            <Label htmlFor="phone">
              Phone Number <span className="text-gray-400">(optional)</span>
            </Label>
            <Input
              id="phone"
              type="tel"
              placeholder="+91 9876543210"
              value={phone}
              onChange={(e) => handlePhoneChange(e.target.value)}
              disabled={isLoading}
            />
            {errors.phone && (
              <p className="text-sm text-red-600">{errors.phone}</p>
            )}
            <p className="text-xs text-gray-500">
              Include country code (e.g., +91 for India, +1 for US)
            </p>
          </div>
        )}

        <Button
          type="submit"
          className="w-full"
          disabled={isLoading || !isFormValid}
        >
          {isLoading ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Creating account...
            </>
          ) : (
            buttonText
          )}
        </Button>
      </form>

      <div className="text-xs text-gray-500 text-center">
        By creating an account, you agree to our{" "}
        <a href="/terms" className="text-blue-600 hover:text-blue-500">
          Terms of Service
        </a>{" "}
        and{" "}
        <a href="/privacy" className="text-blue-600 hover:text-blue-500">
          Privacy Policy
        </a>
      </div>
    </div>
  );
}
