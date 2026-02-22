import { createFileRoute, redirect } from "@tanstack/react-router";
import { useState } from "react";
import { useAuth } from "@/hooks/useAuth";
import { useCategories, useCreateCategory, useUpdateCategory, useDeleteCategory } from "@/hooks/useRisks";
import Loader from "@/components/loader";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { PlusIcon, PencilIcon, TrashIcon, MoreHorizontalIcon } from "lucide-react";
import { toast } from "sonner";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

export const Route = createFileRoute("/app/categories")({
  loader: ({ context }) => {
    const { queryClient } = context;
    const authData = queryClient.getQueryData<{ role: string }>(["auth", "me"]);
    if (!authData || authData.role !== "admin") {
      throw redirect({ to: "/app" });
    }
  },
  component: Categories,
});

function Categories() {
  const { user, isLoading: authLoading } = useAuth();
  const { data: categories, isLoading: categoriesLoading, error } = useCategories();
  const createCategory = useCreateCategory();
  const updateCategory = useUpdateCategory();
  const deleteCategory = useDeleteCategory();

  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [formData, setFormData] = useState({ name: "", description: "" });

  if (authLoading) {
    return (
      <div className="flex items-center justify-center p-8">
        <Loader />
      </div>
    );
  }

  if (user?.role !== "admin") {
    return null;
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name.trim()) {
      toast.error("Category name is required");
      return;
    }

    try {
      if (editingId) {
        await updateCategory.mutateAsync({ id: editingId, ...formData });
        toast.success("Category updated successfully");
        setEditingId(null);
        setIsFormOpen(false);
      } else {
        await createCategory.mutateAsync(formData);
        toast.success("Category created successfully");
        setIsFormOpen(false);
      }
      setFormData({ name: "", description: "" });
    } catch (err: any) {
      toast.error(err.message || "Failed to save category");
    }
  };

  const handleEdit = (category: any) => {
    setEditingId(category.id);
    setFormData({ name: category.name, description: category.description || "" });
    setIsFormOpen(true);
  };

  const handleDelete = async (id: string) => {
    if (window.confirm("Are you sure you want to delete this category?")) {
      try {
        await deleteCategory.mutateAsync(id);
        toast.success("Category deleted successfully");
      } catch (err: any) {
        toast.error(err.message || "Failed to delete category");
      }
    }
  };

  const handleCancel = () => {
    setIsFormOpen(false);
    setEditingId(null);
    setFormData({ name: "", description: "" });
  };

  const isSubmitting = createCategory.isPending || updateCategory.isPending;

  return (
    <div className="p-8">
      <div className="flex items-center justify-end mb-6">
        <Button onClick={() => {
            setEditingId(null);
            setFormData({ name: "", description: "" });
            setIsFormOpen(true);
        }}>
          <PlusIcon className="size-4 mr-2" />
          Add Category
        </Button>
      </div>

      {isFormOpen && (
        <Card className="mb-6">
          <CardHeader>
            <CardTitle>{editingId ? "Edit Category" : "Create New Category"}</CardTitle>
            <CardDescription>{editingId ? "Update category details" : "Add a new risk category"}</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <Input
                  placeholder="Category name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                />
                <Input
                  placeholder="Description (optional)"
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                />
              </div>
              <div className="flex gap-2">
                <Button type="submit" disabled={isSubmitting}>
                  {isSubmitting ? "Saving..." : (editingId ? "Update" : "Create")}
                </Button>
                <Button type="button" variant="outline" onClick={handleCancel}>
                  Cancel
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      {categoriesLoading ? (
        <div className="flex items-center justify-center p-8">
          <Loader />
        </div>
      ) : error ? (
        <div className="text-center p-8 text-muted-foreground">
          Failed to load categories. You may not have admin access.
        </div>
      ) : categories?.length === 0 ? (
        <div className="text-center p-8 text-muted-foreground">
          No categories found. Create one to get started.
        </div>
      ) : (
        <div className="grid gap-4">
          {categories?.map((category) => (
            <Card key={category.id}>
              <CardHeader className="flex flex-row items-start justify-between space-y-0">
                <div className="space-y-1">
                    <CardTitle className="text-lg">{category.name}</CardTitle>
                    {category.description && (
                    <CardDescription>{category.description}</CardDescription>
                    )}
                </div>
                <DropdownMenu>
                  <DropdownMenuTrigger>
                    <Button variant="ghost" size="icon" className="-mr-2 h-8 w-8">
                      <MoreHorizontalIcon className="h-4 w-4" />
                      <span className="sr-only">Actions</span>
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem onClick={() => handleEdit(category)}>
                      <PencilIcon className="mr-2 h-4 w-4" />
                      Edit
                    </DropdownMenuItem>
                    <DropdownMenuItem 
                        onClick={() => handleDelete(category.id)}
                        className="text-destructive focus:text-destructive"
                    >
                      <TrashIcon className="mr-2 h-4 w-4" />
                      Delete
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </CardHeader>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
