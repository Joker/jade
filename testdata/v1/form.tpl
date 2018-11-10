<form method="post">
    <fieldset>
        <legend>General</legend>
        <p>
            <label for="user[name]">Username:
                <input type="text" name="user[name]" value="{{ print user.name }}"/>
            </label>
        </p>
        <p>
            <label for="user[email]">Email:
                <input type="text" name="user[email]" value="{{ print user.email }}"/>
                <div class="tip">
                    Enter a valid 
                    email address 
                    such as <em>tj@vision-media.ca</em>.
                </div>
            </label>
        </p>
    </fieldset>
    <fieldset>
        <legend>Location</legend>
        <p>
            <label for="user[city]">City:
                <input type="text" name="user[city]" value="{{ print user.city }}"/>
            </label>
        </p>
        <p>
            <select name="user[province]">
                <option value="">-- Select Province --</option>
                <option value="AB">Alberta</option>
                <option value="BC">British Columbia</option>
                <option value="SK">Saskatchewan</option>
                <option value="MB">Manitoba</option>
                <option value="ON">Ontario</option>
                <option value="QC">Quebec</option>
            </select>
        </p>
    </fieldset>
    <p class="buttons">
        <input type="submit" value="Save"/>
    </p>
</form>
