class CreateItem < ActiveRecord::Migration[5.2]
  def change
    create_table :items do |t|
      t.string :name, null: false

      # Article to use when referring to the item.
      # Due to proper and improper variations, the
      # thing field is case sensitive.
      # options: "the", "some", "a", "an", "The"
      #
      t.string :article, null: false

      # effects
      t.integer :xp_incr,        null: false, default: 0
      t.integer :level_incr,     null: false, default: 0
      t.integer :strength_incr,  null: false, default: 0
      t.integer :charisma_incr,  null: false, default: 0
      t.integer :intellect_incr, null: false, default: 0
      t.integer :agility_incr,   null: false, default: 0
      t.integer :luck_incr,      null: false, default: 0
    end
  end
end
