package commands

import (
	"fmt"
	"strconv"

	"github.com/fintrack/fintrack/internal/db"
	"github.com/fintrack/fintrack/internal/db/repositories"
	"github.com/fintrack/fintrack/internal/models"
	"github.com/fintrack/fintrack/internal/output"
	"github.com/spf13/cobra"
)

// NewCategoryCmd creates the category command
func NewCategoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "category",
		Aliases: []string{"cat", "c"},
		Short:   "Manage transaction categories",
		Long: `Manage transaction categories for organizing income and expenses.
Categories can be hierarchical with parent-child relationships.`,
	}

	cmd.AddCommand(newCategoryAddCmd())
	cmd.AddCommand(newCategoryListCmd())
	cmd.AddCommand(newCategoryShowCmd())
	cmd.AddCommand(newCategoryUpdateCmd())
	cmd.AddCommand(newCategoryDeleteCmd())

	return cmd
}

func newCategoryAddCmd() *cobra.Command {
	var parentName string
	var color string
	var icon string

	cmd := &cobra.Command{
		Use:   "add <name> <type>",
		Short: "Add a new category",
		Long: `Add a new transaction category.

Type must be one of: income, expense, transfer

Examples:
  fintrack category add "Groceries" expense
  fintrack category add "Salary" income --color "#00FF00"
  fintrack cat add "Coffee" expense --parent "Food & Dining"`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			categoryType := args[1]

			// Validate category type
			if categoryType != models.CategoryTypeIncome &&
				categoryType != models.CategoryTypeExpense &&
				categoryType != models.CategoryTypeTransfer {
				return fmt.Errorf("invalid category type: %s (must be income, expense, or transfer)", categoryType)
			}

			repo := repositories.NewCategoryRepository(db.Get())

			// Check if category already exists
			exists, err := repo.NameExists(name, categoryType, nil)
			if err != nil {
				return fmt.Errorf("failed to check category existence: %w", err)
			}
			if exists {
				return fmt.Errorf("category '%s' of type '%s' already exists", name, categoryType)
			}

			category := &models.Category{
				Name:     name,
				Type:     categoryType,
				Color:    color,
				Icon:     icon,
				IsSystem: false,
			}

			// Handle parent category
			if parentName != "" {
				parent, err := repo.GetByName(parentName, categoryType)
				if err != nil {
					return fmt.Errorf("parent category not found: %w", err)
				}
				category.ParentID = &parent.ID
			}

			if err := repo.Create(category); err != nil {
				return output.PrintError(cmd, fmt.Errorf("failed to create category: %w", err))
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				return output.Print(cmd, category)
			}

			fmt.Printf("Category created successfully (ID: %d)\n", category.ID)
			fmt.Printf("Name: %s\n", category.Name)
			fmt.Printf("Type: %s\n", category.Type)
			if category.ParentID != nil {
				fmt.Printf("Parent: %s\n", parentName)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&parentName, "parent", "", "Parent category name")
	cmd.Flags().StringVar(&color, "color", "", "Color code (e.g., #FF5733)")
	cmd.Flags().StringVar(&icon, "icon", "", "Icon identifier")

	return cmd
}

func newCategoryListCmd() *cobra.Command {
	var categoryType string
	var topLevelOnly bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all categories",
		Long: `List all transaction categories.

Examples:
  fintrack category list
  fintrack category list --type expense
  fintrack cat list --top-level`,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo := repositories.NewCategoryRepository(db.Get())

			var categories []*models.Category
			var err error

			if topLevelOnly {
				categories, err = repo.ListTopLevel(categoryType)
			} else {
				categories, err = repo.List(categoryType)
			}

			if err != nil {
				return output.PrintError(cmd, fmt.Errorf("failed to list categories: %w", err))
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				return output.Print(cmd, categories)
			}

			if len(categories) == 0 {
				fmt.Println("No categories found.")
				return nil
			}

			// Create table
			table := output.NewTable("ID", "Name", "Type", "Parent", "System", "Color", "Icon")
			for _, cat := range categories {
				parentName := ""
				if cat.Parent != nil {
					parentName = cat.Parent.Name
				}

				systemStr := "No"
				if cat.IsSystem {
					systemStr = "Yes"
				}

				table.AddRow(
					strconv.FormatUint(uint64(cat.ID), 10),
					cat.Name,
					cat.Type,
					parentName,
					systemStr,
					cat.Color,
					cat.Icon,
				)
			}

			table.Print()
			fmt.Printf("\nTotal: %d categories\n", len(categories))

			return nil
		},
	}

	cmd.Flags().StringVarP(&categoryType, "type", "t", "", "Filter by type (income, expense, transfer)")
	cmd.Flags().BoolVar(&topLevelOnly, "top-level", false, "Show only top-level categories")

	return cmd
}

func newCategoryShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <id>",
		Short: "Show category details",
		Long: `Display detailed information about a specific category.

Examples:
  fintrack category show 5
  fintrack cat show 10`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid category ID: %w", err)
			}

			repo := repositories.NewCategoryRepository(db.Get())
			category, err := repo.GetByID(uint(id))
			if err != nil {
				return output.PrintError(cmd, err)
			}

			// Get subcategories
			subcategories, err := repo.ListSubcategories(category.ID)
			if err != nil {
				return output.PrintError(cmd, fmt.Errorf("failed to get subcategories: %w", err))
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				data := map[string]interface{}{
					"category":      category,
					"subcategories": subcategories,
				}
				return output.Print(cmd, data)
			}

			fmt.Printf("Category ID: %d\n", category.ID)
			fmt.Printf("Name: %s\n", category.Name)
			fmt.Printf("Type: %s\n", category.Type)
			if category.Parent != nil {
				fmt.Printf("Parent: %s (ID: %d)\n", category.Parent.Name, category.Parent.ID)
			} else {
				fmt.Printf("Parent: (none)\n")
			}
			fmt.Printf("System Category: %v\n", category.IsSystem)
			if category.Color != "" {
				fmt.Printf("Color: %s\n", category.Color)
			}
			if category.Icon != "" {
				fmt.Printf("Icon: %s\n", category.Icon)
			}
			fmt.Printf("Created: %s\n", category.CreatedAt.Format("2006-01-02 15:04:05"))

			if len(subcategories) > 0 {
				fmt.Printf("\nSubcategories (%d):\n", len(subcategories))
				for _, sub := range subcategories {
					fmt.Printf("  - %s (ID: %d)\n", sub.Name, sub.ID)
				}
			}

			return nil
		},
	}

	return cmd
}

func newCategoryUpdateCmd() *cobra.Command {
	var name string
	var color string
	var icon string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a category",
		Long: `Update an existing category's properties.

Examples:
  fintrack category update 5 --name "New Name"
  fintrack cat update 10 --color "#FF0000" --icon "🏠"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid category ID: %w", err)
			}

			repo := repositories.NewCategoryRepository(db.Get())
			category, err := repo.GetByID(uint(id))
			if err != nil {
				return output.PrintError(cmd, err)
			}

			// Check if it's a system category
			if category.IsSystem {
				return fmt.Errorf("cannot update system category")
			}

			// Apply updates
			updated := false
			if cmd.Flags().Changed("name") {
				// Check for duplicate names
				exists, err := repo.NameExists(name, category.Type, &category.ID)
				if err != nil {
					return fmt.Errorf("failed to check category existence: %w", err)
				}
				if exists {
					return fmt.Errorf("category '%s' of type '%s' already exists", name, category.Type)
				}
				category.Name = name
				updated = true
			}

			if cmd.Flags().Changed("color") {
				category.Color = color
				updated = true
			}

			if cmd.Flags().Changed("icon") {
				category.Icon = icon
				updated = true
			}

			if !updated {
				return fmt.Errorf("no updates specified")
			}

			if err := repo.Update(category); err != nil {
				return output.PrintError(cmd, fmt.Errorf("failed to update category: %w", err))
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				return output.Print(cmd, category)
			}

			fmt.Println("Category updated successfully")
			fmt.Printf("Name: %s\n", category.Name)
			fmt.Printf("Type: %s\n", category.Type)

			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "New category name")
	cmd.Flags().StringVar(&color, "color", "", "New color code")
	cmd.Flags().StringVar(&icon, "icon", "", "New icon identifier")

	return cmd
}

func newCategoryDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a category",
		Long: `Delete a category. System categories cannot be deleted.

Examples:
  fintrack category delete 5
  fintrack cat delete 10`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid category ID: %w", err)
			}

			repo := repositories.NewCategoryRepository(db.Get())

			if err := repo.Delete(uint(id)); err != nil {
				return output.PrintError(cmd, fmt.Errorf("failed to delete category: %w", err))
			}

			if output.GetFormat(cmd) == output.FormatJSON {
				return output.Print(cmd, map[string]string{"message": "Category deleted successfully"})
			}

			fmt.Println("Category deleted successfully")
			return nil
		},
	}

	return cmd
}
