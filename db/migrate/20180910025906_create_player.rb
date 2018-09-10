class CreatePlayer < ActiveRecord::Migration[5.2]
  def change
    create_table :players do |t|
      t.string :handle, null: false
      t.string :name, null: false
      t.string :class, null: false
      t.integer :xp, null: false, default: 0
      t.timestamps
    end
  end
end
